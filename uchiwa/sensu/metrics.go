package sensu

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sensu/uchiwa/uchiwa/structs"
)

// Metric returns the Sensu Enterprise metrics for the clients
func (s *Sensu) Metric(name string) (*structs.SERawMetric, error) {
	if name == "" {
		return nil, errors.New("Metric name can't be empty")
	}

	body, _, err := s.getBytes(fmt.Sprintf("%s/%s", "metrics", name))
	if err != nil {
		return nil, err
	}

	var metric structs.SERawMetric
	if err := json.Unmarshal(body, &metric); err != nil {
		return nil, fmt.Errorf("Parsing JSON-encoded response body: %v", err)
	}

	return &metric, nil
}
