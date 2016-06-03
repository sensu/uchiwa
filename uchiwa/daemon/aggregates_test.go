package daemon

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestBuildAggregates(t *testing.T) {
	data := structs.Data{
		Aggregates: []interface{}{
			map[string]interface{}{"dc": "us-east-1", "name": "foo"},
			map[string]interface{}{"dc": "us-east-1", "check": "bar"},
		},
	}
	d := Daemon{Data: &data}
	d.buildAggregates()

	aggregate1 := d.Data.Aggregates[0].(map[string]interface{})
	aggregate2 := d.Data.Aggregates[1].(map[string]interface{})

	assert.Equal(t, aggregate1["name"], "foo")
	assert.Nil(t, aggregate1["check"])
	assert.Equal(t, aggregate2["name"], "bar")
	assert.Nil(t, aggregate2["check"])
}
