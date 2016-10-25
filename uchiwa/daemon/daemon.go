package daemon

import (
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
	setID(d.Data.Stashes, ":/")
	d.BuildSubscriptions()
	setID(d.Data.Aggregates, "/")
	d.buildMetrics()
	d.buildSEMetrics()
}

// getData retrieves all endpoints for every datacenter
func (d *Daemon) fetchData() {
	d.Data.Health.Sensu = make(map[string]structs.SensuHealth, len(*d.Datacenters))

	for _, datacenter := range *d.Datacenters {
		logger.Infof("Updating the datacenter %s", datacenter.Name)

		// set default health status
		d.Data.Health.Sensu[datacenter.Name] = structs.SensuHealth{Output: datacenterErrorString, Status: 2}
		d.Data.Health.Uchiwa = "ok"

		// fetch sensu data from the datacenter
		stashes, err := datacenter.GetStashes()
		if err != nil {
			logger.Warningf("Connection failed to the datacenter %s", datacenter.Name)
			continue
		}
		silenced, err := datacenter.GetSilenced()
		if err != nil {
			logger.Warningf("Connection failed to the datacenter %s.", datacenter.Name)
			continue
		}
		checks, err := datacenter.GetChecks()
		if err != nil {
			logger.Warningf("Connection failed to the datacenter %s", datacenter.Name)
			continue
		}
		clients, err := datacenter.GetClients()
		if err != nil {
			logger.Warningf("Connection failed to the datacenter %s", datacenter.Name)
			continue
		}
		events, err := datacenter.GetEvents()
		if err != nil {
			logger.Warningf("Connection failed to the datacenter %s", datacenter.Name)
			continue
		}
		info, err := datacenter.GetInfo()
		if err != nil {
			logger.Warningf("Connection failed to the datacenter %s", datacenter.Name)
			continue
		}
		aggregates, err := datacenter.GetAggregates()
		if err != nil {
			logger.Warningf("Connection failed to the datacenter %s", datacenter.Name)
			continue
		}

		if d.Enterprise {
			d.Data.SERawMetrics = *getEnterpriseMetrics(&datacenter, &d.Data.SERawMetrics)
		}

		// Determine the status of the datacenter
		if !info.Redis.Connected {
			d.Data.Health.Sensu[datacenter.Name] = structs.SensuHealth{Output: "Not connected to Redis", Status: 1}
		} else if !info.Transport.Connected {
			d.Data.Health.Sensu[datacenter.Name] = structs.SensuHealth{Output: "Not connected to the transport", Status: 1}
		} else {
			d.Data.Health.Sensu[datacenter.Name] = structs.SensuHealth{Output: "ok", Status: 0}
		}

		// add fetched data into d.Data interface
		for _, v := range stashes {
			setDc(v, datacenter.Name)
			d.Data.Stashes = append(d.Data.Stashes, v)
		}
		for _, v := range silenced {
			setDc(v, datacenter.Name)
			d.Data.Silenced = append(d.Data.Silenced, v)
		}
		for _, v := range checks {
			setDc(v, datacenter.Name)
			d.Data.Checks = append(d.Data.Checks, v)
		}
		for _, v := range clients {
			setDc(v, datacenter.Name)
			d.Data.Clients = append(d.Data.Clients, v)
		}
		for _, v := range events {
			setDc(v, datacenter.Name)
			d.Data.Events = append(d.Data.Events, v)
		}
		for _, v := range aggregates {
			setDc(v, datacenter.Name)
			d.Data.Aggregates = append(d.Data.Aggregates, v)
		}

		// build datacenter
		dc := d.buildDatacenter(&datacenter.Name, info)
		dc.Stats["aggregates"] = len(aggregates)
		dc.Stats["checks"] = len(checks)
		dc.Stats["clients"] = len(clients)
		dc.Stats["events"] = len(events)
		dc.Stats["silenced"] = len(silenced)
		dc.Stats["stashes"] = len(stashes)
		d.Data.Dc = append(d.Data.Dc, dc)
	}
}

func (d *Daemon) resetData() {
	d.Data = &structs.Data{}
}

// getEnterpriseMetrics retrieves Sensu Enterprise metrics
func getEnterpriseMetrics(datacenter SensuDatacenter, metrics *structs.SERawMetrics) *structs.SERawMetrics {
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
	metrics.Clients = append(metrics.Clients, m["clients"])
	metrics.Events = append(metrics.Events, m["events"])
	metrics.KeepalivesAVG60 = append(metrics.KeepalivesAVG60, m["keepalives_avg_60"])
	metrics.Requests = append(metrics.Requests, m["check_requests"])
	metrics.Results = append(metrics.Results, m["results"])

	return metrics
}
