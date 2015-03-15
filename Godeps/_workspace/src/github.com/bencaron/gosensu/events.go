package sensu

import (
	"encoding/json"
	"fmt"
)

// GetEvents Return all the current events
func (s *Sensu) GetEvents() ([]interface{}, error) {
	return s.getList("events", 0, 0)
}

// GetEventsForClient Returns the current events for given client
func (s *Sensu) GetEventsForClient(client string) ([]interface{}, error) {
	//return s.get("events", client)
	// TODO  is this the correct way? need validation??
	return s.getList(fmt.Sprintf("events/%s", client), 0, 0)
}

// GetEventsCheckForClient Returns the event for a check for a client
func (s *Sensu) GetEventsCheckForClient(client string, check string) ([]interface{}, error) {
	//return s.get("events", client)
	// TODO  is this the correct way? need validation??
	return s.getList(fmt.Sprintf("events/%s/%s", client, check), 0, 0)
}

// ResolveEvent delete an event
func (s *Sensu) ResolveEvent(payload interface{}) (map[string]interface{}, error) {
	//	return s.post(fmt.Sprintf("stashes/create"), payload)
	payloadstr, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Stash parsing error: %q returned: %v", err, err)
	}
	return s.postPayload("resolve", string(payloadstr[:]))
}
