package uchiwa

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestFindCheck(t *testing.T) {
	u := Uchiwa{
		Data: &structs.Data{},
	}

	u.Data.Checks = []interface{}{
		map[string]interface{}{"name": "foo", "dc": "us-east-1"},
		map[string]interface{}{"name": "bar", "dc": "us-east-1"},
		map[string]interface{}{"name": "foo", "dc": "us-west-1"},
	}

	checks, err := u.findCheck("foo")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(checks))

	checks, err = u.findCheck("bar")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(checks))

	_, err = u.findCheck("qux")
	assert.NotNil(t, err)
}
