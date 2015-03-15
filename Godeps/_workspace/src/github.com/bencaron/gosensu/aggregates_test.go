package sensu

import (
	"fmt"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetAggregates(t *testing.T) {
	assert := assert.New(t)
	sensu := getSensuTester()

	if strings.Contains(sensu.URL, "localhost") {
		t.Skip("Not testing aggregates with canned")
	}
	if assert.NotNil(t, sensu) {
		agg, err := sensu.GetAggregates()
		assert.Nil(err, fmt.Sprintf("GetAggregates return an error: %s", err))
		assert.NotNil(agg, "GetAggregates return a nil result!")
	}
}
