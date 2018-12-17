package daemon

import (
	"context"
	"fmt"
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
	Info       *structs.Info
	Silenced   []interface{}
	Stashes    []interface{}
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

// Fetch retrieves all data for a given datacenter
func (f *DatacenterFetcher) Fetch() {
	defer f.wg.Done()

	logger.Infof("updating the datacenter %s", f.datacenter.Name)

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

	start := time.Now()
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	// ctx.Done() := make(chan bool, 1)

	// fetch sensu data from the datacenter
	wg.Add(7)
	go d.fetchStashes(ctx, errCh)
	go d.fetchSilenced(ctx, errCh)
	go d.fetchChecks(ctx, errCh)
	go d.fetchClients(ctx, errCh)
	go d.fetchEvents(ctx, errCh)
	go d.fetchInfo(ctx, errCh)
	go d.fetchAggregates(ctx, errCh)

	if f.enterprise {
		wg.Add(1)
		go d.fetchEnterpriseMetrics(ctx, errCh)
	}

	// Wait for all goroutines to complete. Close the ctx.Done() once all
	// goroutines have properly returned.
	go func() {
		wg.Wait()
		cancel()
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		if err != nil {
			// Stop all goroutines
			cancel()

			// Log the error
			logger.Warning(err)
			elapsed := time.Since(start)
			logger.Warningf("failed to update the datacenter %s in %s", d.datacenter.Name, elapsed)

			// Mark the datacenter as down
			f.mutex.Lock()
			f.data.Health.Sensu[f.datacenter.Name] = structs.SensuHealth{
				Output: err.Error(), Status: 2,
			}
			f.mutex.Unlock()
			return
		}
	}

	// update health
	f.mutex.Lock()
	f.data.Health.Sensu[f.datacenter.Name] = d.determineHealth()
	f.mutex.Unlock()

	// build datacenter
	dc := f.buildDatacenter(&d.datacenter.Name, d.snapshot.Info)
	dc.Metrics["aggregates"] = len(d.snapshot.Aggregates)
	dc.Metrics["checks"] = len(d.snapshot.Checks)
	dc.Metrics["clients"] = len(d.snapshot.Clients)
	dc.Metrics["events"] = len(d.snapshot.Events)
	dc.Metrics["silenced"] = len(d.snapshot.Silenced)
	dc.Metrics["stashes"] = len(d.snapshot.Stashes)

	// update datacenter in the Daemon scope
	f.mutex.Lock()
	f.data.Dc = append(f.data.Dc, dc)
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

	elapsed := time.Since(start)
	logger.Infof("updated the datacenter %s in %s", f.datacenter.Name, elapsed)
}

func (d *DatacenterSnapshotFetcher) determineHealth() structs.SensuHealth {
	if d.snapshot.Info != nil {
		if !d.snapshot.Info.Redis.Connected {
			return structs.SensuHealth{Output: "Not connected to Redis", Status: 1}
		}
		if !d.snapshot.Info.Transport.Connected && d.snapshot.Info != nil {
			return structs.SensuHealth{Output: "Not connected to the transport", Status: 1}
		}
	}

	return structs.SensuHealth{Output: "ok", Status: 0}
}

func (d *DatacenterSnapshotFetcher) fetchStashes(ctx context.Context, errCh chan error) {
	defer d.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			stashes, err := d.datacenter.GetStashes(ctx)
			if err != nil {
				errCh <- fmt.Errorf(
					"could not retrieve stashes from datacenter %s: %s",
					d.datacenter.Name, err)
			}

			d.mutex.Lock()
			for _, v := range stashes {
				setDc(v, d.datacenter.Name)
				d.snapshot.Stashes = append(d.snapshot.Stashes, v)
			}
			d.mutex.Unlock()
			return
		}
	}
}

func (d *DatacenterSnapshotFetcher) fetchSilenced(ctx context.Context, errCh chan error) {
	defer d.wg.Done()

	for {
		select {
		case _ = <-ctx.Done():
			return
		default:
			silenced, err := d.datacenter.GetSilenced(ctx)
			if err != nil {
				errCh <- fmt.Errorf(
					"could not retrieve silenced entries from datacenter %s: %s",
					d.datacenter.Name, err)
			}
			d.mutex.Lock()
			for _, v := range silenced {
				setDc(v, d.datacenter.Name)
				d.snapshot.Silenced = append(d.snapshot.Silenced, v)
			}
			d.mutex.Unlock()
			return
		}
	}
}

func (d *DatacenterSnapshotFetcher) fetchChecks(ctx context.Context, errCh chan error) {
	defer d.wg.Done()

	for {
		select {
		case _ = <-ctx.Done():
			return
		default:
			checks, err := d.datacenter.GetChecks(ctx)
			if err != nil {
				errCh <- fmt.Errorf(
					"could not retrieve checks from datacenter %s: %s",
					d.datacenter.Name, err)
			}
			d.mutex.Lock()

			for _, v := range checks {
				setDc(v, d.datacenter.Name)
				d.snapshot.Checks = append(d.snapshot.Checks, v)
			}

			d.mutex.Unlock()
			return
		}
	}

}

func (d *DatacenterSnapshotFetcher) fetchClients(ctx context.Context, errCh chan error) {
	defer d.wg.Done()

	for {
		select {
		case _ = <-ctx.Done():
			return
		default:
			clients, err := d.datacenter.GetClients(ctx)
			if err != nil {
				errCh <- fmt.Errorf(
					"could not retrieve clients entries from datacenter %s: %s",
					d.datacenter.Name, err)
			}
			d.mutex.Lock()

			for _, v := range clients {
				setDc(v, d.datacenter.Name)
				d.snapshot.Clients = append(d.snapshot.Clients, v)
			}

			d.mutex.Unlock()
			return
		}
	}

}

func (d *DatacenterSnapshotFetcher) fetchEvents(ctx context.Context, errCh chan error) {
	defer d.wg.Done()

	for {
		select {
		case _ = <-ctx.Done():
			return
		default:
			events, err := d.datacenter.GetEvents(ctx)
			if err != nil {
				errCh <- fmt.Errorf(
					"could not retrieve events from datacenter %s: %s",
					d.datacenter.Name, err)
			}
			d.mutex.Lock()
			for _, v := range events {
				setDc(v, d.datacenter.Name)
				d.snapshot.Events = append(d.snapshot.Events, v)
			}

			d.mutex.Unlock()
			return
		}
	}
}

func (d *DatacenterSnapshotFetcher) fetchInfo(ctx context.Context, errCh chan error) {
	defer d.wg.Done()

	for {
		select {
		case _ = <-ctx.Done():
			return
		default:
			info, err := d.datacenter.GetInfo()
			if err != nil {
				errCh <- fmt.Errorf(
					"could not retrieve info about datacenter %s: %s",
					d.datacenter.Name, err)
			}
			d.mutex.Lock()
			d.snapshot.Info = info
			d.mutex.Unlock()
			return
		}
	}
}

func (d *DatacenterSnapshotFetcher) fetchAggregates(ctx context.Context, errCh chan error) {
	defer d.wg.Done()

	for {
		select {
		case _ = <-ctx.Done():
			logger.Warning("stopping aggregates")
			return
		default:
			aggregates, err := d.datacenter.GetAggregates(ctx)
			if err != nil {
				errCh <- fmt.Errorf(
					"could not retrieve aggregates from datacenter %s: %s",
					d.datacenter.Name, err)
			}
			d.mutex.Lock()

			for _, v := range aggregates {
				setDc(v, d.datacenter.Name)
				d.snapshot.Aggregates = append(d.snapshot.Aggregates, v)
			}

			d.mutex.Unlock()
			return
		}
	}
}

func (d *DatacenterSnapshotFetcher) fetchEnterpriseMetrics(ctx context.Context, errCh chan error) {
	defer d.wg.Done()

	d.mutex.Lock()
	d.metrics = getEnterpriseMetrics(&d.datacenter)
	d.mutex.Unlock()
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
