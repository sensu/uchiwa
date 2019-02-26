package sensu

import (
	"encoding/json"
	"fmt"

	"github.com/sensu/uchiwa/uchiwa/structs"
)

// GetInfo returns a pointer to a structs.Info struct containing the
// Sensu version and the transport and Redis connection information
func (s *Sensu) GetInfo() (*structs.Info, error) {
	body, _, err := s.getBytes("info")
	if err != nil {
		return nil, err
	}

	var info structs.Info
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("Parsing JSON-encoded response body: %v", err)
	}

	return &info, nil
}

// GetInfo returns a pointer to a structs.Info struct containing the
// Sensu version and the transport and Redis connection information
func (s *Sensu) GetInfoFromAPI(i int) (*structs.Info, error) {
	api := &s.APIs[i]
	body, _, err := s.getBytesFromAPI(api, "info")
	if err != nil {
		return nil, err
	}

	var info structs.Info
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("Parsing JSON-encoded response body: %v", err)
	}

	return &info, nil
}
