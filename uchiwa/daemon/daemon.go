package daemon

import (
	"time"

	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/uchiwa/sensu"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

const datacenterErrorString = "Connection error. Is the Sensu API running?"

// Daemon structure is used to manage the Uchiwa daemon
type Daemon struct {
	Data        *structs.Data
	Datacenters *[]sensu.Sensu
}

// Start method fetches and builds Sensu data from each datacenter every Refresh seconds
func (d *Daemon) Start(interval int, data chan *structs.Data) {
	// immediately fetch the first set of data and send it over the data channel
	d.fetchData()
	d.buildData()

	select {
	case data <- d.Data:
		logger.Debug("Sending initial results on the 'data' channel")
	default:
		logger.Debug("Could not send initial results on the 'data' channel")
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
			logger.Debug("Sending results on the 'data' channel")
		default:
			logger.Debug("Could not send results on the 'data' channel")
		}
	}
}

// buildData method prepares fetched data
func (d *Daemon) buildData() {
	d.buildEvents()
	d.buildClients()
	d.BuildSubscriptions()
	d.buildResults()
}

// getData retrieves all endpoints for every datacenter
func (d *Daemon) fetchData() {
	d.Data.Health.Sensu = make(map[string]structs.SensuHealth, len(*d.Datacenters))

	for _, datacenter := range *d.Datacenters {
		// set default health status
		d.Data.Health.Sensu[datacenter.Name] = structs.SensuHealth{Output: datacenterErrorString}
		d.Data.Health.Uchiwa = "ok"

		// fetch sensu data from the datacenter
		stashes, err := datacenter.GetStashes()
		if err != nil {
			logger.Warning(err)
			continue
		}
		checks, err := datacenter.GetChecks()
		if err != nil {
			logger.Warning(err)
			continue
		}
		clients, err := datacenter.GetClients()
		if err != nil {
			logger.Warning(err)
			continue
		}
		events, err := datacenter.GetEvents()
		if err != nil {
			logger.Warning(err)
			continue
		}
		info, err := datacenter.Info()
		if err != nil {
			logger.Warning(err)
			continue
		}
		aggregates, err := datacenter.GetAggregates()
		if err != nil {
			logger.Warning(err)
			continue
		}
		results, err := datacenter.Results()
		if err == nil {
			for _, v := range *results {
				setDc(v, datacenter.Name)
				d.Data.Results = append(d.Data.Results, v)
			}
		}

		d.Data.Health.Sensu[datacenter.Name] = structs.SensuHealth{Output: "ok"}

		// add fetched data into d.Data interface
		for _, v := range stashes {
			setDc(v, datacenter.Name)
			d.Data.Stashes = append(d.Data.Stashes, v)
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
		dc.Stats["stashes"] = len(stashes)
		d.Data.Dc = append(d.Data.Dc, dc)
	}
}

func (d *Daemon) resetData() {
	d.Data = &structs.Data{}
}
