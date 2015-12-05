package daemon

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestBuildChecks(t *testing.T) {
	data := structs.Data{
		Checks: []interface{}{
			map[string]interface{}{"dc": "us-east-1", "name": "foo"},
			map[string]interface{}{"dc": "us-east-1", "name": "bar"},
			map[string]interface{}{"dc": "us-west-1", "name": "foo"},
		},
	}
	d := Daemon{Data: &data}
	d.buildChecks()

	check1 := d.Data.Checks[0].(map[string]interface{})
	check2 := d.Data.Checks[1].(map[string]interface{})
	check3 := d.Data.Checks[2].(map[string]interface{})

	assert.Equal(t, check1["_id"], "us-east-1/foo")
	assert.Equal(t, check2["_id"], "us-east-1/bar")
	assert.Equal(t, check3["_id"], "us-west-1/foo")
}
