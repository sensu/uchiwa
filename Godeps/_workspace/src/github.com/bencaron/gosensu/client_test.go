package sensu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClientList(t *testing.T) {
	assert := assert.New(t)
	sensu := getSensuTester()
	events, err := sensu.GetEvents()
	assert.Nil(err, fmt.Sprintf("GetClientList returned an error: %s", err))
	assert.NotNil(events, "GetClientList returned nil!")
}
