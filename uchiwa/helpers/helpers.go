package helpers

import (
	"net"
	"net/http"

	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// BuildClientsMetrics builds the metrics for the events
func BuildClientsMetrics(clients *[]interface{}) *structs.StatusMetrics {
	metrics := structs.StatusMetrics{}

	metrics.Total = len(*clients)

	for _, c := range *clients {
		client := c.(map[string]interface{})

		status, ok := client["status"].(int)
		if !ok {
			logger.Warningf("Could not assert this status to an int: %+v", client["status"])
			continue
		}

		if status == 2.0 {
			metrics.Critical++
			continue
		} else if status == 1.0 {
			metrics.Warning++
			continue
		} else if status == 0.0 {
			continue
		}
		metrics.Unknown++
	}

	return &metrics
}

// BuildEventsMetrics builds the metrics for the events
func BuildEventsMetrics(events *[]interface{}) *structs.StatusMetrics {
	metrics := structs.StatusMetrics{}

	metrics.Total = len(*events)

	for _, e := range *events {
		event := e.(map[string]interface{})

		check, ok := event["check"].(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert this check to an interface: %+v", event["check"])
			continue
		}

		status, ok := check["status"].(float64)
		if !ok {
			logger.Warningf("Could not assert this status to a flot64: %+v", check["status"])
			continue
		}

		if status == 2.0 {
			metrics.Critical++
			continue
		} else if status == 1.0 {
			metrics.Warning++
			continue
		}
		metrics.Unknown++
	}

	return &metrics
}

// GetIP returns the real user IP address
func GetIP(r *http.Request) string {
	if xForwardedFor := r.Header.Get("X-FORWARDED-FOR"); len(xForwardedFor) > 0 {
		return xForwardedFor
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
