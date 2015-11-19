package uchiwa

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/config"
	"github.com/stretchr/testify/assert"
)

func TestInitDatacenters(t *testing.T) {
	// A single datacenter
	conf := config.Config{
		Sensu: []config.SensuConfig{
			{Name: "foo", URL: "http://10.0.0.1:4567"},
		},
	}
	datacenters := initDatacenters(&conf)
	assert.Equal(t, 1, len(*datacenters))
	assert.Equal(t, 1, len((*datacenters)[0].APIs))
	assert.Equal(t, "foo", (*datacenters)[0].Name)

	// Two datacenters
	conf = config.Config{
		Sensu: []config.SensuConfig{
			{Name: "foo", URL: "http://10.0.0.1:4567"},
			{Name: "bar", URL: "http://10.0.0.2:4567"},
		},
	}
	datacenters = initDatacenters(&conf)
	assert.Equal(t, 2, len(*datacenters))
	assert.Equal(t, 1, len((*datacenters)[1].APIs))
	assert.Equal(t, "bar", (*datacenters)[1].Name)

	// One datacenter with three APIs
	conf = config.Config{
		Sensu: []config.SensuConfig{
			{Name: "foo", URL: "http://10.0.0.1:4567"},
			{Name: "foo", URL: "http://10.0.0.2:4567"},
			{Name: "foo", URL: "http://10.0.0.3:4567"},
		},
	}
	datacenters = initDatacenters(&conf)
	assert.Equal(t, 1, len(*datacenters))
	assert.Equal(t, 3, len((*datacenters)[0].APIs))
	assert.Equal(t, "foo", (*datacenters)[0].Name)

	// Two datacenters with four APIs
	conf = config.Config{
		Sensu: []config.SensuConfig{
			{Name: "foo", URL: "http://10.0.0.1:4567"},
			{Name: "bar", URL: "http://10.0.0.10:4567"},
			{Name: "foo", URL: "http://10.0.0.1:4567"},
			{Name: "bar", URL: "http://10.0.0.11:4567"},
		},
	}
	datacenters = initDatacenters(&conf)
	assert.Equal(t, 2, len(*datacenters))
	assert.Equal(t, 2, len((*datacenters)[0].APIs))
	assert.Equal(t, 2, len((*datacenters)[1].APIs))
	assert.Equal(t, "foo", (*datacenters)[0].Name)
	assert.Equal(t, "bar", (*datacenters)[1].Name)
}
