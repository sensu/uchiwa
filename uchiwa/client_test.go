package uchiwa

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestBuildClientHistory(t *testing.T) {
	var dc string
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

	client := map[string]interface{}{"name": "qux"}
	dc = "us-east-1"
	status = Success
	history = []interface{}{map[string]interface{}{"check": "cpu", "last_result": map[string]interface{}{"command": "cpu.rb", "status": status}, "last_status": status}}
	expectedHistory = []interface{}{map[string]interface{}{"check": "cpu", "client": "qux", "dc": "us-east-1", "last_result": map[string]interface{}{"command": "cpu.rb", "status": status}, "last_status": status, "silenced": false, "silenced_by": []string(nil)}}
	result := u.buildClientHistory(client, dc, history)
	assert.Equal(t, expectedHistory, result)
}

func TestFindClient(t *testing.T) {
	u := Uchiwa{
		Data: &structs.Data{},
	}

	u.Data.Clients = []interface{}{
		map[string]interface{}{"name": "foo", "dc": "us-east-1"},
		map[string]interface{}{"name": "bar", "dc": "us-east-1"},
		map[string]interface{}{"name": "foo", "dc": "us-west-1"},
	}

	clients, err := u.findClient("foo")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(clients))

	clients, err = u.findClient("bar")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(clients))

	_, err = u.findClient("qux")
	assert.NotNil(t, err)
}
