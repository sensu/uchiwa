package sensu

import "fmt"

// GetEvents returns a slice of all clients
func (s *Sensu) GetEvents() ([]interface{}, error) {
	return s.getSlice("events", DefaultLimit)
}

// DeleteEvent delete an event
func (s *Sensu) DeleteEvent(check, client string) error {
	return s.delete(fmt.Sprintf("events/%s/%s", client, check))
}
