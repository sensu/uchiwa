package uchiwa

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestBuildClientHistory(t *testing.T) {
	var client, dc string
	var status float64
	const Critical float64 = 2.0
	const Warning float64 = 1.0
	const Success float64 = 0.0
	var history = []interface{}{}
	var expectedHistory = []interface{}{}

	u := Uchiwa{
		Data: &structs.Data{},
	}
	u.Data.Events = []interface{}{
		map[string]interface{}{"action": "create", "check": map[string]interface{}{"name": "cpu", "command": "cpu.rb", "status": Critical}, "client": map[string]interface{}{"name": "foo"}, "dc": "us-east-1", "occurrences": 7},
		map[string]interface{}{"check": map[string]interface{}{"name": "cpu", "command": "cpu.rb", "status": Warning}, "client": map[string]interface{}{"name": "bar"}, "dc": "us-east-1", "occurrences": 5},
		map[string]interface{}{"check": "cpu", "client": "qux", "dc": "us-west-1", "occurrences": 10, "output": "CRITICAL", "status": Critical},
	}
	u.Data.Stashes = []interface{}{map[string]interface{}{"dc": "us-east-1", "path": "silence/foo/cpu"}}

	// Sensu => 0.18; we already have the last_result attribute in history
	client = "foo"
	dc = "us-east-1"
	status = Critical
	history = []interface{}{map[string]interface{}{"check": "cpu", "last_result": map[string]interface{}{"command": "cpu.rb", "status": status}, "last_status": status}}
	expectedHistory = []interface{}{map[string]interface{}{"acknowledged": true, "check": "cpu", "client": "foo", "dc": "us-east-1", "last_result": map[string]interface{}{"action": "create", "command": "cpu.rb", "name": "cpu", "occurrences": 7, "status": status}, "last_status": status}}
	result := u.buildClientHistory(client, dc, history)
	assert.Equal(t, expectedHistory, result)

	client = "qux"
	dc = "us-east-1"
	status = Success
	history = []interface{}{map[string]interface{}{"check": "cpu", "last_result": map[string]interface{}{"command": "cpu.rb", "status": status}, "last_status": status}}
	expectedHistory = []interface{}{map[string]interface{}{"acknowledged": false, "check": "cpu", "client": "qux", "dc": "us-east-1", "last_result": map[string]interface{}{"command": "cpu.rb", "status": status}, "last_status": status}}
	result = u.buildClientHistory(client, dc, history)
	assert.Equal(t, expectedHistory, result)

	// 0.12 > Sensu < 0.18; we don't have the last_result attribute in history but we have rich events
	client = "bar"
	dc = "us-east-1"
	status = Warning
	history = []interface{}{map[string]interface{}{"check": "cpu", "last_status": status}}
	expectedHistory = []interface{}{map[string]interface{}{"acknowledged": false, "check": "cpu", "client": "bar", "dc": "us-east-1", "last_result": map[string]interface{}{"command": "cpu.rb", "name": "cpu", "occurrences": 5, "status": status}, "last_status": status}}
	result = u.buildClientHistory(client, dc, history)
	assert.Equal(t, expectedHistory, result)

	// Sensu <= 0.12; we don't have the last_result attribute in history and no rich events
	client = "qux"
	dc = "us-west-1"
	status = Critical
	history = []interface{}{map[string]interface{}{"check": "cpu", "last_status": status}}
	expectedHistory = []interface{}{map[string]interface{}{"acknowledged": false, "check": "cpu", "client": "qux", "dc": "us-west-1", "last_result": map[string]interface{}{"check": "cpu", "client": "qux", "occurrences": 10, "output": "CRITICAL", "status": status}, "last_status": status}}
	result = u.buildClientHistory(client, dc, history)
	assert.Equal(t, expectedHistory, result)

	client = "baz"
	dc = "us-west-1"
	status = Success
	history = []interface{}{map[string]interface{}{"check": "cpu", "last_execution": 1445195709, "last_status": status}}
	expectedHistory = []interface{}{map[string]interface{}{"acknowledged": false, "check": "cpu", "client": "baz", "dc": "us-west-1", "last_execution": 1445195709, "last_result": map[string]interface{}{"last_execution": 1445195709, "status": status}, "last_status": status}}
	result = u.buildClientHistory(client, dc, history)
	assert.Equal(t, expectedHistory, result)

}
