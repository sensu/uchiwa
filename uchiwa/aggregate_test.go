package uchiwa

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestFindAggregate(t *testing.T) {
	u := Uchiwa{
		Data: &structs.Data{},
	}

	u.Data.Aggregates = []interface{}{
		map[string]interface{}{"name": "foo", "dc": "us-east-1"},
		map[string]interface{}{"name": "bar", "dc": "us-east-1"},
		map[string]interface{}{"name": "foo", "dc": "us-west-1"},
	}

	aggregates, err := u.findAggregate("foo")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(aggregates))

	aggregates, err = u.findAggregate("bar")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(aggregates))

	_, err = u.findAggregate("qux")
	assert.NotNil(t, err)
}
