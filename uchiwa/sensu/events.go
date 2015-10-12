package sensu

import "fmt"

// GetEvents Return all the current events
func (s *Sensu) GetEvents() ([]interface{}, error) {
	return s.getList("events", NoLimit)
}

// GetEventsForClient Returns the current events for given client
func (s *Sensu) GetEventsForClient(client string) ([]interface{}, error) {
	return s.getList(fmt.Sprintf("events/%s", client), NoLimit)
}

// GetEventsCheckForClient Returns the event for a check for a client
func (s *Sensu) GetEventsCheckForClient(client string, check string) ([]interface{}, error) {
	return s.getList(fmt.Sprintf("events/%s/%s", client, check), NoLimit)
}

// ResolveEvent delete an event
func (s *Sensu) ResolveEvent(check, client string) error {
	return s.delete(fmt.Sprintf("events/%s/%s", client, check))
}
