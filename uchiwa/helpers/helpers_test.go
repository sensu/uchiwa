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

func TestIsAcknowledged(t *testing.T) {
	var check, client, dc string
	var stashes = []interface{}{}

	isAcknowledged := IsAcknowledged(check, client, dc, stashes)
	assert.Equal(t, false, isAcknowledged, "the dc string is empty")

	dc = "us-east-1"
	isAcknowledged = IsAcknowledged(check, client, dc, stashes)
	assert.Equal(t, false, isAcknowledged, "the client string is empty")

	client = "foo"
	isAcknowledged = IsAcknowledged(check, client, dc, stashes)
	assert.Equal(t, false, isAcknowledged, "the stashes slice is empty")

	stashes = []interface{}{map[string]interface{}{"dc": "us-west-1", "path": "silence/foo"}}
	isAcknowledged = IsAcknowledged(check, client, dc, stashes)
	assert.Equal(t, false, isAcknowledged, "the client's datacenter was not properly compared")

	stashes = []interface{}{map[string]interface{}{"dc": "us-east-1", "path": "silence/foo/bar"}}
	isAcknowledged = IsAcknowledged(check, client, dc, stashes)
	assert.Equal(t, false, isAcknowledged, "the client's check was not properly compared")

	check = "bar"
	isAcknowledged = IsAcknowledged(check, client, dc, stashes)
	assert.Equal(t, true, isAcknowledged)

	stashes = []interface{}{map[string]interface{}{"dc": "us-west-1", "path": "silence/foo/bar"}, map[string]interface{}{"dc": "us-east-1", "path": "silence/foo/bar"}}
	isAcknowledged = IsAcknowledged(check, client, dc, stashes)
	assert.Equal(t, true, isAcknowledged)

	stashes = []interface{}{"not_a_map", map[string]interface{}{"dc": "us-east-1", "path": "silence/foo/bar"}}
	isAcknowledged = IsAcknowledged(check, client, dc, stashes)
	assert.Equal(t, true, isAcknowledged)
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
