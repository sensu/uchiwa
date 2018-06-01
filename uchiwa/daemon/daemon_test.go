package daemon

import (
	"errors"
	"testing"

	"github.com/sensu/uchiwa/uchiwa/structs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSensu struct {
	mock.Mock
}

func (s *mockSensu) Metric(name string) (*structs.SERawMetric, error) {
	args := s.Called(name)
	return args.Get(0).(*structs.SERawMetric), args.Error(1)
}
func (s *mockSensu) GetName() string {
	args := s.Called()
	return args.String(0)
}

func TestGetEnterpriseMetrics(t *testing.T) {
	clients := structs.SERawMetric{
		Points: [][]interface{}{[]interface{}{1453578490, 1}, []interface{}{1453578500, 1}},
	}

	datacenter := new(mockSensu)
	datacenter.On("GetName").Return("foo")
	datacenter.On("Metric", "clients").Return(&clients, nil)
	datacenter.On("Metric", "events").Return(&structs.SERawMetric{}, errors.New(""))
	datacenter.On("Metric", "keepalives_avg_60").Return(&structs.SERawMetric{}, nil)
	datacenter.On("Metric", "check_requests").Return(&structs.SERawMetric{}, nil)
	datacenter.On("Metric", "results").Return(&structs.SERawMetric{}, nil)
	metrics := structs.SERawMetrics{}

	metrics = getEnterpriseMetrics(datacenter)

	assert.Equal(t, 2, len(metrics.Clients[0].Points))
	assert.Equal(t, 0, len(metrics.Events[0].Points))
	assert.Equal(t, 0, len(metrics.KeepalivesAVG60[0].Points))

	datacenter.AssertExpectations(t)
}
