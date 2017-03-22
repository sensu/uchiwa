package sensu

import (
	"encoding/json"
	"fmt"
)

// DeleteClient deletes a client using its name
func (s *Sensu) DeleteClient(client string) error {
	return s.delete(fmt.Sprintf("clients/%s", client))
}

// GetClients returns a slice of all clients
func (s *Sensu) GetClients() ([]interface{}, error) {
	return s.getSlice("clients", DefaultLimit)
}

// GetClient returns a map of a specific client corresponding to the provided client name
func (s *Sensu) GetClient(client string) (map[string]interface{}, error) {
	return s.getMap(fmt.Sprintf("clients/%s", client))
}

// GetClientHistory returns a slice containing the history of a specific check corresponding to the provided client name
func (s *Sensu) GetClientHistory(client string) ([]interface{}, error) {
	return s.getSlice(fmt.Sprintf("clients/%s/history", client), NoLimit)
}

func (s *Sensu) UpdateClient(payload interface{}) (map[string]interface{}, error) {
	payloadstr, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing the payload: %s", err)
	}

	return s.postPayload(fmt.Sprintf("clients"), string(payloadstr[:]))
}
