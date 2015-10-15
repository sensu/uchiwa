package sensu

import "fmt"

// GetClients Return the list of clients
func (s *Sensu) GetClients() ([]interface{}, error) {
	return s.getList("clients", DefaultLimit)
}

// GetClient Return client info
func (s *Sensu) GetClient(client string) (map[string]interface{}, error) {
	return s.getMap(fmt.Sprintf("clients/%s", client))
}

// GetClientHistory Return client history
func (s *Sensu) GetClientHistory(client string) ([]interface{}, error) {
	return s.getList(fmt.Sprintf("clients/%s/history", client), NoLimit)
}

// DeleteClient Return the list of clients
func (s *Sensu) DeleteClient(client string) error {
	return s.delete(fmt.Sprintf("clients/%s", client))
}
