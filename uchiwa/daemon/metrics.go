package daemon

import (
	"github.com/sensu/uchiwa/uchiwa/helpers"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// buildMetrics ...
func (d *Daemon) buildMetrics() {
	d.Data.Metrics.Aggregates.Total = len(d.Data.Aggregates)
	d.Data.Metrics.Checks.Total = len(d.Data.Checks)
	d.Data.Metrics.Datacenters.Total = len(d.Data.Dc)
	d.Data.Metrics.Silenced.Total = len(d.Data.Silenced)
	d.Data.Metrics.Stashes.Total = len(d.Data.Stashes)

	d.Data.Metrics.Clients = *helpers.BuildClientsMetrics(&d.Data.Clients)
	d.Data.Metrics.Events = *helpers.BuildEventsMetrics(&d.Data.Events)
}

// buildSEMetrics prepares the Sensu Enterprise metrics for frontend consumption by the HUD
func (d *Daemon) buildSEMetrics() {
	d.Data.SEMetrics.Clients = rawMetricsToAggregatedCoordinates(d.Data.SERawMetrics.Clients)
	if len(d.Data.SEMetrics.Clients.Data) > 361 {
		d.Data.SEMetrics.Clients.Data = d.Data.SEMetrics.Clients.Data[len(d.Data.SEMetrics.Clients.Data)-361:]
	}
	d.Data.SEMetrics.Clients.Name = "Clients"

	d.Data.SEMetrics.Events = rawMetricsToCoordinates(d.Data.SERawMetrics.Events)

	d.Data.SEMetrics.KeepalivesAVG60 = rawMetricsToAggregatedCoordinates(d.Data.SERawMetrics.KeepalivesAVG60)
	if len(d.Data.SEMetrics.KeepalivesAVG60.Data) > 361 {
		d.Data.SEMetrics.KeepalivesAVG60.Data = d.Data.SEMetrics.KeepalivesAVG60.Data[len(d.Data.SEMetrics.KeepalivesAVG60.Data)-361:]
	}
	d.Data.SEMetrics.KeepalivesAVG60.Name = "Keepalives"

	d.Data.SEMetrics.Requests = rawMetricsToAggregatedCoordinates(d.Data.SERawMetrics.Requests)
	d.Data.SEMetrics.Requests.Name = "Requests"

	d.Data.SEMetrics.Results = rawMetricsToAggregatedCoordinates(d.Data.SERawMetrics.Results)
	d.Data.SEMetrics.Results.Name = "Results"
}

// rawMetricsToCoordinates converts raw metrics from Sensu API to frontend-ready data
// E.g. the points ([]float32{123456789, 0.5}) are transformed to coordinates ([]{X: 123456789000, Y: 0.5})
func rawMetricsToCoordinates(rawMetrics []*structs.SERawMetric) []*structs.SEMetric {
	seMetrics := make([]*structs.SEMetric, len(rawMetrics))

	for i, metrics := range rawMetrics {
		total := len(metrics.Points)
		data := make([]structs.XY, total)
		total--
		for j, point := range metrics.Points {
			if len(point) != 2 {
				continue
			}

			data[j] = structs.XY{X: point[0].(float64) * 1000, Y: point[1].(float64)}
		}
		if len(data) > 120 {
			data = data[len(data)-120:]
		}
		seMetrics[i] = &structs.SEMetric{
			Data: data,
			Name: metrics.Name,
		}
	}
	return seMetrics
}

// rawMetricsToAggregatedCoordinates converts raw metrics from Sensu API to frontend-ready data
// E.g. the slice of slices ([][]float32) aggregates all slices and is transformed to coordinates ([]{X: 123456789000, Y: 0.5})
func rawMetricsToAggregatedCoordinates(rawMetrics []*structs.SERawMetric) *structs.SEMetric {
	// Find the oldest data point in the last position
	var oldest float64
	for _, metrics := range rawMetrics {
		count := len(metrics.Points)
		if count == 0 || len(metrics.Points[count-1]) == 0 {
			metrics.Points = make([][]interface{}, 0)
			break
		}

		timestamp, ok := metrics.Points[count-1][0].(float64)
		if !ok {
			metrics.Points = make([][]interface{}, 0)
			break
		}

		if oldest == 0 {
			oldest = timestamp
		} else if oldest > timestamp {
			oldest = timestamp
		}
	}

	// Make sure all arrays have the same latest point
	for _, metrics := range rawMetrics {
		count := len(metrics.Points)
		if count == 0 || len(metrics.Points[count-1]) == 0 {
			metrics.Points = make([][]interface{}, 0)
			break
		}

		value, ok := metrics.Points[count-1][0].(float64)
		if !ok {
			metrics.Points = make([][]interface{}, 0)
			break
		}
		if value == oldest {
			continue
		} else {
			// This metrics contains a too recent data point that we need to remove
			_, metrics.Points = metrics.Points[len(metrics.Points)-1], metrics.Points[:len(metrics.Points)-1]
		}
	}

	// Get the length of the biggest array
	var count int
	for _, metrics := range rawMetrics {
		if len(metrics.Points) > count {
			count = len(metrics.Points)
		}
	}

	var seMetric structs.SEMetric
	seMetric.Data = make([]structs.XY, count)
	for i := range seMetric.Data {
		for _, metrics := range rawMetrics {
			if len(metrics.Points) > i {
				if len(metrics.Points[len(metrics.Points)-i-1]) < 2 {
					continue
				}

				x, ok := metrics.Points[len(metrics.Points)-i-1][0].(float64)
				if !ok {
					continue
				}

				y, ok := metrics.Points[len(metrics.Points)-i-1][1].(float64)
				if !ok {
					continue
				}

				if seMetric.Data[count-i-1].X == 0 {
					seMetric.Data[count-i-1].X = x * 1000
					seMetric.Data[count-i-1].Y = y
				} else if seMetric.Data[count-i-1].X == (x * 1000) {
					seMetric.Data[count-i-1].Y += y
				}
			}
		}
	}
	return &seMetric
}
