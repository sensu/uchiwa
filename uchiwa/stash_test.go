package uchiwa

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestFindStash(t *testing.T) {
	u := Uchiwa{
		Data: &structs.Data{},
	}

	u.Data.Stashes = []interface{}{
		map[string]interface{}{"path": "foo", "dc": "us-east-1"},
		map[string]interface{}{"path": "silence/foo", "dc": "us-east-1"},
		map[string]interface{}{"path": "bar", "dc": "us-east-1"},
		map[string]interface{}{"path": "foo", "dc": "us-west-1"},
	}

	stashes, err := u.findStash("foo")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(stashes))

	stashes, err = u.findStash("silence/foo")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(stashes))

	stashes, err = u.findStash("bar")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(stashes))

	_, err = u.findCheck("qux")
	assert.NotNil(t, err)
}
