package sensu

import (
	"encoding/json"
	"fmt"

	"github.com/sensu/uchiwa/uchiwa/structs"
)

// Info Will return the Sensu version along with rabbitmq and redis information.
func (s *Sensu) Info() (*structs.Info, error) {
	body, _, err := s.get(fmt.Sprintf("%s/%s", s.URL, "info"))
	if err != nil {
		return nil, err
	}

	var info structs.Info
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("Parsing JSON-encoded response body: %v", err)
	}

	return &info, nil
}
