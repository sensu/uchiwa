package sensu

import "fmt"

// GetClients Return the list of clients
func (s *Sensu) GetClients() ([]interface{}, error) {
	return s.getList("clients", 0, 0)
}

// GetClientsSlice Return a slice in the list of clients
func (s *Sensu) GetClientsSlice(limit int, offset int) ([]interface{}, error) {
	return s.getList("clients", limit, offset)
}

// GetClient Return client info
func (s *Sensu) GetClient(client string) (map[string]interface{}, error) {
	return s.getMap(fmt.Sprintf("clients/%s", client))
}

// GetClientHistory Return client history
func (s *Sensu) GetClientHistory(client string) ([]interface{}, error) {
	return s.getList(fmt.Sprintf("clients/%s/history", client), 0, 0)
}

// DeleteClient Return the list of clients
func (s *Sensu) DeleteClient(client string) error {
	return s.delete(fmt.Sprintf("clients/%s", client))
}
