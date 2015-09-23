package helpers

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestBuildClientsMetrics(t *testing.T) {
	clients := []interface{}{map[string]interface{}{"status": 0}, map[string]interface{}{"status": 1}, map[string]interface{}{"status": 2}, map[string]interface{}{"status": 3}}
	expectedMetrics := structs.StatusMetrics{Critical: 1, Total: 4, Unknown: 1, Warning: 1}

	metrics := BuildClientsMetrics(&clients)
	assert.Equal(t, expectedMetrics, *metrics)
}

func TestBuildEventsMetrics(t *testing.T) {
	events := []interface{}{map[string]interface{}{"check": map[string]interface{}{"status": 1.0}}, map[string]interface{}{"check": map[string]interface{}{"status": 2.0}}, map[string]interface{}{"check": map[string]interface{}{"status": 3.0}}}
	expectedMetrics := structs.StatusMetrics{Critical: 1, Total: 3, Unknown: 1, Warning: 1}

	metrics := BuildEventsMetrics(&events)
	assert.Equal(t, expectedMetrics, *metrics)
}
