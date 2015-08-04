package daemon

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestBuildResult(t *testing.T) {
	d := Daemon{}

	d.Data = &structs.Data{
		Results: []interface{}{
			map[string]interface{}{"client": "vodka", "check": map[string]interface{}{"name": "hydrogen"}, "dc": "foo"},
			map[string]interface{}{"client": "whisky", "check": map[string]interface{}{"name": "helium"}, "dc": "foo"},
			map[string]interface{}{"client": "vodka", "check": map[string]interface{}{"name": "hydrogen"}, "dc": "bar"},
		},
		Stashes: []interface{}{map[string]interface{}{"dc": "foo", "path": "silence/vodka/hydrogen"}},
	}

	d.buildResults()

	expectedResults := []interface{}{
		map[string]interface{}{"client": "vodka", "check": map[string]interface{}{"acknowledged": true, "name": "hydrogen"}, "dc": "foo"},
		map[string]interface{}{"client": "whisky", "check": map[string]interface{}{"acknowledged": false, "name": "helium"}, "dc": "foo"},
		map[string]interface{}{"client": "vodka", "check": map[string]interface{}{"acknowledged": false, "name": "hydrogen"}, "dc": "bar"},
	}
	assert.Equal(t, expectedResults, d.Data.Results)
}
