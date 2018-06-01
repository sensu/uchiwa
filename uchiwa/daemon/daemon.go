package daemon

import (
	"sync"
	"time"

	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/sensu"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

const datacenterErrorString = "Connection error. Is the Sensu API running?"

// Daemon structure is used to manage the Uchiwa daemon
type Daemon struct {
	Data        *structs.Data
	Datacenters *[]sensu.Sensu
	Enterprise  bool
}

// DatacenterFetcher is used to manage the fetching of data from a datacenter
type DatacenterFetcher struct {
	data       *structs.Data
	datacenter sensu.Sensu
	mutex      *sync.Mutex
	wg         *sync.WaitGroup
	enterprise bool
}

// DatacenterSnapshotFetcher is used to manage the fetching of data from a datacenter API endpoint
type DatacenterSnapshotFetcher struct {
	snapshot   *DatacenterSnapshot
	metrics    structs.SERawMetrics
	datacenter sensu.Sensu
	mutex      *sync.Mutex
	wg         *sync.WaitGroup
}

// DatacenterSnapshot is used to store a snapshot of a datacenter's data
type DatacenterSnapshot struct {
	Aggregates []interface{}
	Checks     []interface{}
	Clients    []interface{}
	Events     []interface{}
	Info       structs.Info
	Silenced   []interface{}
	Stashes    []interface{}
	Error      string
}

// SensuDatacenter represents the sensu.Sensu struct
type SensuDatacenter interface {
	GetName() string
	Metric(string) (*structs.SERawMetric, error)
}

// Start method fetches and builds Sensu data from each datacenter every Refresh seconds
func (d *Daemon) Start(interval int, data chan *structs.Data) {
	// immediately fetch the first set of data and send it over the data channel
	d.fetchData()
	d.buildData()

	select {
	case data <- d.Data:
		logger.Trace("Sending initial results on the 'data' channel")
	default:
		logger.Trace("Could not send initial results on the 'data' channel")
	}

	// fetch new data every interval
	duration := time.Duration(interval) * time.Second
	for _ = range time.Tick(duration) {
		d.resetData()
		d.fetchData()
		d.buildData()

		// send the result over the data channel
		select {
		case data <- d.Data:
			logger.Trace("Sending results on the 'data' channel")
		default:
			logger.Trace("Could not send results on the 'data' channel")
		}
	}
}

// buildData method prepares fetched data
func (d *Daemon) buildData() {
	d.buildEvents()
	d.buildClients()
	setID(d.Data.Checks, "/")
	setID(d.Data.Silenced, ":")
	setID(d.Data.Stashes, "/")
	d.BuildSubscriptions()
	setID(d.Data.Aggregates, "/")
	d.buildMetrics()
	d.buildSEMetrics()
}

// fetchData retrieves all data from each datacenter
func (d *Daemon) fetchData() {
	d.Data.Health.Sensu = make(map[string]structs.SensuHealth, len(*d.Datacenters))

	mutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	for _, datacenter := range *d.Datacenters {
		dc := DatacenterFetcher{
			data:       d.Data,
			datacenter: datacenter,
			mutex:      mutex,
			wg:         wg,
			enterprise: d.Enterprise,
		}

		wg.Add(1)
		go dc.Fetch()
	}

	wg.Wait()
}

// fetch retrieves all data for a given datacenter
func (f *DatacenterFetcher) Fetch() {
	logger.Infof("Updating the datacenter %s", f.datacenter.Name)

	// set default health status
	f.mutex.Lock()
	f.data.Health.Sensu[f.datacenter.Name] = structs.SensuHealth{Output: datacenterErrorString, Status: 2}
	f.data.Health.Uchiwa = "ok"
	f.mutex.Unlock()

	mutex := &sync.Mutex{}
	wg := &sync.WaitGroup{}

	d := DatacenterSnapshotFetcher{
		snapshot:   &DatacenterSnapshot{},
		datacenter: f.datacenter,
		mutex:      mutex,
		wg:         wg,
	}

	// fetch sensu data from the datacenter
	wg.Add(8)
	go d.fetchStashes()
	go d.fetchSilenced()
	go d.fetchChecks()
	go d.fetchClients()
	go d.fetchEvents()
	go d.fetchInfo()
	go d.fetchAggregates()

	if f.enterprise {
		go d.fetchEnterpriseMetrics()
	}

	wg.Wait()

	// build datacenter
	dc := f.buildDatacenter(&d.datacenter.Name, &d.snapshot.Info)
	dc.Metrics["aggregates"] = len(d.snapshot.Aggregates)
	dc.Metrics["checks"] = len(d.snapshot.Checks)
	dc.Metrics["clients"] = len(d.snapshot.Clients)
	dc.Metrics["events"] = len(d.snapshot.Events)
	dc.Metrics["silenced"] = len(d.snapshot.Silenced)
	dc.Metrics["stashes"] = len(d.snapshot.Stashes)

	// update datacenter in the Daemon scope
	f.mutex.Lock()
	f.data.Dc = append(f.data.Dc, dc)
	f.data.Health.Sensu[f.datacenter.Name] = d.determineHealth()
	f.data.Stashes = append(f.data.Stashes, d.snapshot.Stashes...)
	f.data.Silenced = append(f.data.Silenced, d.snapshot.Silenced...)
	f.data.Checks = append(f.data.Checks, d.snapshot.Checks...)
	f.data.Clients = append(f.data.Clients, d.snapshot.Clients...)
	f.data.Events = append(f.data.Events, d.snapshot.Events...)
	f.data.Aggregates = append(f.data.Aggregates, d.snapshot.Aggregates...)

	if f.enterprise {
		f.data.SERawMetrics.Clients = append(f.data.SERawMetrics.Clients, d.metrics.Clients...)
		f.data.SERawMetrics.Events = append(f.data.SERawMetrics.Events, d.metrics.Events...)
		f.data.SERawMetrics.KeepalivesAVG60 = append(f.data.SERawMetrics.KeepalivesAVG60, d.metrics.KeepalivesAVG60...)
		f.data.SERawMetrics.Requests = append(f.data.SERawMetrics.Requests, d.metrics.Requests...)
		f.data.SERawMetrics.Results = append(f.data.SERawMetrics.Results, d.metrics.Results...)
	}

	f.mutex.Unlock()

	logger.Infof("Updated the datacenter %s", f.datacenter.Name)

	f.wg.Done()
}

func (d *DatacenterSnapshotFetcher) determineHealth() structs.SensuHealth {
	if !d.snapshot.Info.Redis.Connected {
		return structs.SensuHealth{Output: "Not connected to Redis", Status: 1}
	} else if !d.snapshot.Info.Transport.Connected {
		return structs.SensuHealth{Output: "Not connected to the transport", Status: 1}
	} else if d.snapshot.Error != "" {
		return structs.SensuHealth{Output: d.snapshot.Error, Status: 2}
	}

	return structs.SensuHealth{Output: "ok", Status: 0}
}

func (d *DatacenterSnapshotFetcher) fetchStashes() {
	stashes, err := d.datacenter.GetStashes()
	d.mutex.Lock()
	if err != nil {
		logger.Debug(err)
		logger.Warningf("Connection failed to the datacenter %s", d.datacenter.Name)
		d.snapshot.Error = err.Error()
	}

	for _, v := range stashes {
		setDc(v, d.datacenter.Name)
		d.snapshot.Stashes = append(d.snapshot.Stashes, v)
	}

	d.mutex.Unlock()
	d.wg.Done()
}

func (d *DatacenterSnapshotFetcher) fetchSilenced() {
	silenced, err := d.datacenter.GetSilenced()
	d.mutex.Lock()
	if err != nil {
		logger.Debug(err)
		logger.Warningf("Impossible to retrieve silenced entries from the "+
			"datacenter %s. Silencing might not be possible, please update Sensu", d.datacenter.Name)
	}

	for _, v := range silenced {
		setDc(v, d.datacenter.Name)
		d.snapshot.Silenced = append(d.snapshot.Stashes, v)
	}

	d.mutex.Unlock()
	d.wg.Done()
}

func (d *DatacenterSnapshotFetcher) fetchChecks() {
	checks, err := d.datacenter.GetChecks()
	d.mutex.Lock()
	if err != nil {
		logger.Debug(err)
		logger.Warningf("Connection failed to the datacenter %s", d.datacenter.Name)
		d.snapshot.Error = err.Error()
	}

	for _, v := range checks {
		setDc(v, d.datacenter.Name)
		d.snapshot.Checks = append(d.snapshot.Checks, v)
	}

	d.mutex.Unlock()
	d.wg.Done()
}

func (d *DatacenterSnapshotFetcher) fetchClients() {
	clients, err := d.datacenter.GetClients()
	d.mutex.Lock()
	if err != nil {
		logger.Debug(err)
		logger.Warningf("Connection failed to the datacenter %s", d.datacenter.Name)
		d.snapshot.Error = err.Error()
	}

	for _, v := range clients {
		setDc(v, d.datacenter.Name)
		d.snapshot.Clients = append(d.snapshot.Clients, v)
	}

	d.mutex.Unlock()
	d.wg.Done()
}

func (d *DatacenterSnapshotFetcher) fetchEvents() {
	events, err := d.datacenter.GetEvents()
	if err != nil {
		logger.Debug(err)
		logger.Warningf("Connection failed to the datacenter %s", d.datacenter.Name)
		d.mutex.Lock()
		d.snapshot.Error = err.Error()
		d.mutex.Unlock()
	}

	for _, v := range events {
		setDc(v, d.datacenter.Name)
		d.snapshot.Events = append(d.snapshot.Events, v)
	}

	d.wg.Done()
}

func (d *DatacenterSnapshotFetcher) fetchInfo() {
	info, err := d.datacenter.GetInfo()
	d.mutex.Lock()
	if err != nil {
		logger.Debug(err)
		logger.Warningf("Connection failed to the datacenter %s", d.datacenter.Name)
		d.snapshot.Error = err.Error()
	}

	d.snapshot.Info = *info
	d.mutex.Unlock()
	d.wg.Done()
}

func (d *DatacenterSnapshotFetcher) fetchAggregates() {
	aggregates, err := d.datacenter.GetAggregates()
	d.mutex.Lock()
	if err != nil {
		logger.Debug(err)
		logger.Warningf("Connection failed to the datacenter %s", d.datacenter.Name)
		d.snapshot.Error = err.Error()
	}

	for _, v := range aggregates {
		setDc(v, d.datacenter.Name)
		d.snapshot.Aggregates = append(d.snapshot.Aggregates, v)
	}

	d.mutex.Unlock()
	d.wg.Done()
}

func (d *DatacenterSnapshotFetcher) fetchEnterpriseMetrics() {
	d.mutex.Lock()
	d.metrics = getEnterpriseMetrics(&d.datacenter)
	d.mutex.Unlock()
	d.wg.Done()
}

func (d *Daemon) resetData() {
	d.Data = &structs.Data{}
}

// getEnterpriseMetrics retrieves Sensu Enterprise metrics
func getEnterpriseMetrics(datacenter SensuDatacenter) structs.SERawMetrics {
	var err error
	m := make(map[string]*structs.SERawMetric)
	metricsEndpoints := []string{"clients", "events", "keepalives_avg_60", "check_requests", "results"}

	for _, metric := range metricsEndpoints {
		m[metric], err = datacenter.Metric(metric)
		if err != nil {
			logger.Debugf("Could not retrieve the %s enterprise metrics. %s", metric, datacenter.GetName())
			m[metric] = &structs.SERawMetric{}
		}
	}

	m["events"].Name = datacenter.GetName()

	metrics := structs.SERawMetrics{}
	metrics.Clients = append(metrics.Clients, m["clients"])
	metrics.Events = append(metrics.Events, m["events"])
	metrics.KeepalivesAVG60 = append(metrics.KeepalivesAVG60, m["keepalives_avg_60"])
	metrics.Requests = append(metrics.Requests, m["check_requests"])
	metrics.Results = append(metrics.Results, m["results"])

	return metrics
}
