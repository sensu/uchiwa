package sensu

import (
	"encoding/json"
	"fmt"
)

// Results returns a list of current check results for all clients
func (s *Sensu) Results() (*[]interface{}, error) {
	body, err := s.get("results")
	if err != nil {
		return nil, err
	}

	var results []interface{}
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("Parsing JSON-encoded response body: %v", err)
	}
	return &results, nil
}
