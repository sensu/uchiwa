package daemon

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
)

func TestRawMetricsToAggregatedCoordinates(t *testing.T) {
	raw1 := structs.SERawMetric{
		Points: [][]interface{}{{1000000000.0, 0.5}, {1000000001.0, 1.0}, {1000000002.0, 1.0}},
	}
	raw2 := structs.SERawMetric{
		Points: [][]interface{}{{1000000000.0, 2.0}, {1000000001.0, 3.0}, {1000000002.0, 4.0}},
	}
	raw3 := structs.SERawMetric{
		Points: [][]interface{}{{1000000000.0, 1.0}, {1000000001.0, 0.0}, {1000000002.0, 0.0}, {1000000003.0, 2.0}},
	}
	raw4 := structs.SERawMetric{
		Points: [][]interface{}{{1000000002.0, 2.5}},
	}

	metrics := []*structs.SERawMetric{&raw1, &raw2, &raw3, &raw4}

	expectedCoordinates := structs.SEMetric{
		Data: []structs.XY{
			{X: 1000000000000, Y: 3.5},
			{X: 1000000001000, Y: 4},
			{X: 1000000002000, Y: 7.5},
		},
	}

	coordinates := rawMetricsToAggregatedCoordinates(metrics)
	assert.Equal(t, &expectedCoordinates, coordinates)
}
