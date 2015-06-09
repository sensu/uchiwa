package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindClientEvents(t *testing.T) {

	// no events
	client := map[string]interface{}{"dc": "us-east-1", "name": "foo"}
	events := []interface{}{nil}
	expectedClient := map[string]interface{}{"dc": "us-east-1", "name": "foo", "status": 0}
	result := findClientEvents(client, &events)
	assert.Equal(t, expectedClient, result)

	var statusFloat = 1.0
	events = []interface{}{
		map[string]interface{}{
			"check":  map[string]interface{}{"output": "http_critical", "status": 2},
			"client": map[string]interface{}{"name": "foo"},
			"dc":     "us-east-1",
		},
		map[string]interface{}{
			"check":  map[string]interface{}{"output": "http_warning", "status": statusFloat},
			"client": map[string]interface{}{"name": "bar"},
			"dc":     "us-west-1",
		},
	}

	// event where status is a float64
	client = map[string]interface{}{"dc": "us-west-1", "name": "bar"}
	expectedClient = map[string]interface{}{"dc": "us-west-1", "name": "bar", "output": "http_warning", "status": 1}
	result = findClientEvents(client, &events)
	assert.Equal(t, expectedClient, result)

	// event where status is an int
	client = map[string]interface{}{"dc": "us-east-1", "name": "foo"}
	expectedClient = map[string]interface{}{"dc": "us-east-1", "name": "foo", "output": "http_critical", "status": 2}
	result = findClientEvents(client, &events)
	assert.Equal(t, expectedClient, result)

	// client has no events
	client = map[string]interface{}{"dc": "us-east-1", "name": "qux"}
	expectedClient = map[string]interface{}{"dc": "us-east-1", "name": "qux", "status": 0}
	result = findClientEvents(client, &events)
	assert.Equal(t, expectedClient, result)

	// client has multiple events
	newEvents := []interface{}{
		map[string]interface{}{
			"check":  map[string]interface{}{"output": "http_critical", "status": 2},
			"client": map[string]interface{}{"name": "qux"},
			"dc":     "us-east-1",
		},
		map[string]interface{}{
			"check":  map[string]interface{}{"output": "http_warning", "status": 1},
			"client": map[string]interface{}{"name": "qux"},
			"dc":     "us-east-1",
		},
	}
	events = append(events, newEvents[0])
	events = append(events, newEvents[1])

	client = map[string]interface{}{"dc": "us-east-1", "name": "qux"}
	expectedClient = map[string]interface{}{"dc": "us-east-1", "name": "qux", "output": "http_critical and 1 more...", "status": 2}
	result = findClientEvents(client, &events)
	assert.Equal(t, expectedClient, result)

	// event with unknown status
	newEvent := map[string]interface{}{
		"check":  map[string]interface{}{"output": "http_unknown", "status": 3},
		"client": map[string]interface{}{"name": "baz"},
		"dc":     "us-west-1",
	}
	events = append(events, newEvent)

	client = map[string]interface{}{"dc": "us-west-1", "name": "baz"}
	expectedClient = map[string]interface{}{"dc": "us-west-1", "name": "baz", "output": "http_unknown", "status": 3}
	result = findClientEvents(client, &events)
	assert.Equal(t, expectedClient, result)
}
