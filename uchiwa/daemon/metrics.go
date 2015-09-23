package daemon

import "github.com/sensu/uchiwa/uchiwa/helpers"

// buildMetrics ...
func (d *Daemon) buildMetrics() {
	d.Data.Metrics.Aggregates.Total = len(d.Data.Aggregates)
	d.Data.Metrics.Checks.Total = len(d.Data.Checks)
	d.Data.Metrics.Datacenters.Total = len(d.Data.Dc)
	d.Data.Metrics.Stashes.Total = len(d.Data.Stashes)

	d.Data.Metrics.Clients = *helpers.BuildClientsMetrics(&d.Data.Clients)
	d.Data.Metrics.Events = *helpers.BuildEventsMetrics(&d.Data.Events)
}
