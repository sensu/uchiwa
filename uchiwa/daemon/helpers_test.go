package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetID(t *testing.T) {
	elements := []interface{}{
		map[string]interface{}{"dc": "us-east-1", "name": "foo"},
		map[string]interface{}{"dc": "us-east-1", "name": "bar"},
		map[string]interface{}{"dc": "us-west-1", "id": "foo"},
	}

	setID(elements, "/")

	element1 := elements[0].(map[string]interface{})
	element2 := elements[1].(map[string]interface{})
	element3 := elements[2].(map[string]interface{})

	assert.Equal(t, element1["_id"], "us-east-1/foo")
	assert.Equal(t, element2["_id"], "us-east-1/bar")
	assert.Equal(t, element3["_id"], "us-west-1/foo")
}
