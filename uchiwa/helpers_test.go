package uchiwa

import (
	"testing"

	"github.com/bencaron/gosensu"
	"github.com/stretchr/testify/assert"
)

func mockDatacenters() {
	datacenters = make([]sensu.Sensu, 2)
	datacenters[0] = sensu.Sensu{Name: "foo"}
	datacenters[1] = sensu.Sensu{Name: "bar"}
}

func TestGetAPI(t *testing.T) {

	mockDatacenters()

	_, err := getAPI("")
	assert.NotNil(t, err)

	_, err = getAPI("qux")
	assert.NotNil(t, err)

	api, err := getAPI("foo")
	assert.Nil(t, err)
	assert.Equal(t, &datacenters[0], api)

	api, err = getAPI("foo")
	assert.Nil(t, err)
	assert.Equal(t, &datacenters[0], api)
}
