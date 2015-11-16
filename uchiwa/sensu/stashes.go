package sensu

import (
	"encoding/json"
	"fmt"
)

// GetStashes returns a slice of all stashes
func (s *Sensu) GetStashes() ([]interface{}, error) {
	return s.getSlice(fmt.Sprintf("stashes"), NoLimit)
}

// GetStash returns a map of a specific stash corresponding to the provided path
func (s *Sensu) GetStash(path string) (map[string]interface{}, error) {
	return s.getMap(fmt.Sprintf("stashes/%s", path))
}

// CreateStash creates a stash by posting the provided interface as a JSON encoded payload
func (s *Sensu) CreateStash(payload interface{}) (map[string]interface{}, error) {
	payloadstr, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Stash parsing error: %q returned: %v", err, err)
	}
	return s.postPayload(fmt.Sprintf("stashes"), string(payloadstr[:]))
}

// DeleteStash deletes a stash using its path
func (s *Sensu) DeleteStash(path string) error {
	return s.delete(fmt.Sprintf("stashes/%s", path))
}
