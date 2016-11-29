package helpers

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestBuildClientsMetrics(t *testing.T) {
	clients := []interface{}{map[string]interface{}{"status": 0}, map[string]interface{}{"status": 1}, map[string]interface{}{"status": 2}, map[string]interface{}{"status": 3}}
	expectedMetrics := structs.StatusMetrics{Critical: 1, Healthy: 1, Total: 4, Unknown: 1, Warning: 1}
	metrics := BuildClientsMetrics(&clients)
	assert.Equal(t, expectedMetrics, *metrics)

	clients = []interface{}{map[string]interface{}{"status": 1}, map[string]interface{}{"silenced": true, "status": 1}, map[string]interface{}{"silenced": false, "status": 2}}
	expectedMetrics = structs.StatusMetrics{Critical: 1, Silenced: 1, Total: 3, Warning: 1}
	metrics = BuildClientsMetrics(&clients)
	assert.Equal(t, expectedMetrics, *metrics)
}

func TestBuildEventsMetrics(t *testing.T) {
	events := []interface{}{map[string]interface{}{"check": map[string]interface{}{"status": 1.0}}, map[string]interface{}{"check": map[string]interface{}{"status": 2.0}}, map[string]interface{}{"check": map[string]interface{}{"status": 3.0}}}
	expectedMetrics := structs.StatusMetrics{Critical: 1, Total: 3, Unknown: 1, Warning: 1}
	metrics := BuildEventsMetrics(&events)
	assert.Equal(t, expectedMetrics, *metrics)

	events = []interface{}{map[string]interface{}{"check": map[string]interface{}{"status": 1.0}}, map[string]interface{}{"check": map[string]interface{}{"status": 2.0}, "silenced": true}}
	expectedMetrics = structs.StatusMetrics{Silenced: 1, Total: 2, Warning: 1}
	metrics = BuildEventsMetrics(&events)
	assert.Equal(t, expectedMetrics, *metrics)
}

func TestGetBoolFromInterface(t *testing.T) {
	i := map[string]interface{}{"foo": true}

	_, err := GetBoolFromInterface(i)
	assert.NotNil(t, err)

	b, err := GetBoolFromInterface(i["foo"])
	assert.Nil(t, err)
	assert.Equal(t, b, true)
}

func TestGetEvent(t *testing.T) {
	var check, client, dc string
	var events = []interface{}{}

	_, err := GetEvent(check, client, dc, &events)
	assert.NotNil(t, err)

	check = "ram"
	_, err = GetEvent(check, client, dc, &events)
	assert.NotNil(t, err)

	client = "bar"
	_, err = GetEvent(check, client, dc, &events)
	assert.NotNil(t, err)

	dc = "us-west-1"
	_, err = GetEvent(check, client, dc, &events)
	assert.NotNil(t, err)

	events = []interface{}{map[string]interface{}{"check": map[string]interface{}{"name": "cpu", "status": "1"}, "client": map[string]interface{}{"name": "foo"}, "dc": "us-east-1"}}
	_, err = GetEvent(check, client, dc, &events)
	assert.NotNil(t, err, "Wrong datacenter")

	dc = "us-east-1"
	_, err = GetEvent(check, client, dc, &events)
	assert.NotNil(t, err, "Wrong client")

	client = "foo"
	_, err = GetEvent(check, client, dc, &events)
	assert.NotNil(t, err, "Wrong check")

	check = "cpu"
	event, err := GetEvent(check, client, dc, &events)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"name": "cpu", "status": "1"}, event)

	// Sensu <= 0.12 events
	events = []interface{}{map[string]interface{}{"check": "cpu", "client": "foo", "dc": "us-east-1", "occurrences": 10, "output": "CRITICAL", "status": "2"}}
	event, err = GetEvent(check, client, dc, &events)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"check": "cpu", "client": "foo", "occurrences": 10, "output": "CRITICAL", "status": "2"}, event)
}

func TestGetInterfacesFromBytes(t *testing.T) {
	bytes := []byte(`{"foo": "bar"}`)
	_, err := GetInterfacesFromBytes(bytes)
	assert.NotNil(t, err)

	bytes = []byte(`[{"foo": "bar"}, {"baz": "qux"}]`)
	expectedInterfaces := []interface{}{map[string]interface{}{"foo": "bar"}, map[string]interface{}{"baz": "qux"}}
	interfaces, err := GetInterfacesFromBytes(bytes)
	assert.Nil(t, err)
	assert.Equal(t, expectedInterfaces, interfaces)
}

func TestGetMapFromBytes(t *testing.T) {
	bytes := []byte(`[{"foo": "bar"}]`)
	m, err := GetMapFromBytes(bytes)
	assert.NotNil(t, err)

	bytes = []byte(``)
	m, err = GetMapFromBytes(bytes)
	assert.Nil(t, err)

	bytes = []byte(`{"foo": "bar"}`)
	expectedMap := map[string]interface{}{"foo": "bar"}
	m, err = GetMapFromBytes(bytes)
	assert.Nil(t, err)
	assert.Equal(t, expectedMap, m)
}

func TestGetMapFromInterface(t *testing.T) {
	i := map[string]interface{}{"foo": "vodka"}
	m := GetMapFromInterface(i)
	assert.Equal(t, "vodka", m["foo"])
}

func TestIsCheckSilenced(t *testing.T) {
	var check map[string]interface{}
	var client, dc string
	var silenced []interface{}
	var isSilencedBy []string

	isSilenced, _ := IsCheckSilenced(check, client, dc, silenced)
	assert.False(t, isSilenced)

	// Not silenced
	check = map[string]interface{}{"name": "check_cpu", "subscribers": []interface{}{"load-balancer"}}
	client = "foo"
	dc = "us-east-1"
	isSilenced, _ = IsCheckSilenced(check, client, dc, silenced)
	assert.False(t, isSilenced)

	// Wrong datacenter
	silenced = []interface{}{map[string]interface{}{"dc": "us-west-1", "client": "foo", "check": "check_cpu"}}
	isSilenced, _ = IsCheckSilenced(check, client, dc, silenced)
	assert.False(t, isSilenced)

	// Silenced check with check
	// e.g. *:check_cpu
	silenced = []interface{}{map[string]interface{}{"dc": "us-east-1", "id": "*:check_cpu"}}
	isSilenced, isSilencedBy = IsCheckSilenced(check, client, dc, silenced)
	assert.True(t, isSilenced)
	assert.Equal(t, "*:check_cpu", isSilencedBy[0])

	// Silenced check with client subscription
	// e.g. client:foo:*
	silenced = []interface{}{map[string]interface{}{"dc": "us-east-1", "id": "client:foo:*", "subscription": "client:foo"}}
	isSilenced, isSilencedBy = IsCheckSilenced(check, client, dc, silenced)
	assert.True(t, isSilenced)
	assert.Equal(t, "client:foo:*", isSilencedBy[0])

	// Silenced check with client and check subscription
	// e.g. client:foo:check_cpu
	silenced = []interface{}{map[string]interface{}{"dc": "us-east-1", "id": "client:foo:check_cpu", "check": "check_cpu", "subscription": "client:foo"}}
	isSilenced, isSilencedBy = IsCheckSilenced(check, client, dc, silenced)
	assert.True(t, isSilenced)
	assert.Equal(t, "client:foo:check_cpu", isSilencedBy[0])

	// Silenced check with subscription only
	// e.g. load-balancer:*
	silenced = []interface{}{map[string]interface{}{"dc": "us-east-1", "id": "load-balancer:*", "subscription": "load-balancer"}}
	isSilenced, isSilencedBy = IsCheckSilenced(check, client, dc, silenced)
	assert.True(t, isSilenced)
	assert.Equal(t, "load-balancer:*", isSilencedBy[0])

	// Silenced check with *check* and *subscription*
	// e.g. load-balancer:check_cpu
	silenced = []interface{}{map[string]interface{}{"dc": "us-east-1", "id": "load-balancer:check_cpu", "check": "check_cpu", "subscription": "load-balancer"}}
	isSilenced, isSilencedBy = IsCheckSilenced(check, client, dc, silenced)
	assert.True(t, isSilenced)
	assert.Equal(t, "load-balancer:check_cpu", isSilencedBy[0])

	// Silenced check with multiple subscriptions
	silenced = append(silenced, map[string]interface{}{"dc": "us-east-1", "id": "load-balancer:*", "subscription": "load-balancer"})
	silenced = append(silenced, map[string]interface{}{"dc": "us-east-1", "id": "client:foo:*", "subscription": "client:foo"})
	isSilenced, isSilencedBy = IsCheckSilenced(check, client, dc, silenced)
	assert.True(t, isSilenced)
	assert.Equal(t, 3, len(isSilencedBy))
	assert.Equal(t, "load-balancer:check_cpu", isSilencedBy[0])
	assert.Equal(t, "load-balancer:*", isSilencedBy[1])
	assert.Equal(t, "client:foo:*", isSilencedBy[2])

	// Standalone check
	check = map[string]interface{}{"name": "check_cpu"}
	silenced = []interface{}{map[string]interface{}{"dc": "us-east-1", "id": "*:check_cpu"}}
	isSilenced, isSilencedBy = IsCheckSilenced(check, client, dc, silenced)
	assert.True(t, isSilenced)
	assert.Equal(t, "*:check_cpu", isSilencedBy[0])
}

func TestIsClientSilenced(t *testing.T) {
	var client, dc string
	var silenced []interface{}

	isSilenced := IsClientSilenced(client, dc, silenced)
	assert.False(t, isSilenced)

	// Not silenced
	client = "foo"
	dc = "us-east-1"
	isSilenced = IsClientSilenced(client, dc, silenced)
	assert.False(t, isSilenced)

	// Wrong datacenter
	silenced = append(silenced, map[string]interface{}{"dc": "us-west-1", "id": "client:foo:*"})
	isSilenced = IsClientSilenced(client, dc, silenced)
	assert.False(t, isSilenced)

	// Only a check of the client
	silenced = append(silenced, map[string]interface{}{"dc": "us-east-1", "id": "client:foo:check_cpu"})
	isSilenced = IsClientSilenced(client, dc, silenced)
	assert.False(t, isSilenced)

	// Silenced client
	silenced = append(silenced, map[string]interface{}{"dc": "us-east-1", "id": "client:foo:*"})
	isSilenced = IsClientSilenced(client, dc, silenced)
	assert.True(t, isSilenced)
}

func TestIsStringInArray(t *testing.T) {
	var item string
	var array []string

	found := IsStringInArray(item, array)
	assert.Equal(t, false, found, "if item and array are both empty, it should return false")

	item = "foo"
	found = IsStringInArray(item, array)
	assert.Equal(t, false, found, "if array is empty, it should return false")

	array = []string{"bar", "qux"}
	found = IsStringInArray(item, array)
	assert.Equal(t, false, found, "it should return false if the item isn't found in the array")

	array = append(array, "foo")
	found = IsStringInArray(item, array)
	assert.Equal(t, true, found, "it should return true if the item is found in the array")
}
