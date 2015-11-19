package sensu

import "fmt"

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

// DeleteClient deletes a client using its name
func (s *Sensu) DeleteClient(client string) error {
	return s.delete(fmt.Sprintf("clients/%s", client))
}
