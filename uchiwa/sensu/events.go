package sensu

import "fmt"

// GetEvents Return all the current events
func (s *Sensu) GetEvents() ([]interface{}, error) {
	return s.getList("events", 0, 0)
}

// GetEventsForClient Returns the current events for given client
func (s *Sensu) GetEventsForClient(client string) ([]interface{}, error) {
	return s.getList(fmt.Sprintf("events/%s", client), 0, 0)
}

// GetEventsCheckForClient Returns the event for a check for a client
func (s *Sensu) GetEventsCheckForClient(client string, check string) ([]interface{}, error) {
	return s.getList(fmt.Sprintf("events/%s/%s", client, check), 0, 0)
}

// ResolveEvent delete an event
func (s *Sensu) ResolveEvent(check, client string) error {
	return s.delete(fmt.Sprintf("events/%s/%s", client, check))
}
