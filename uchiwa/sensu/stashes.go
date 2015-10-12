package sensu

import (
	"encoding/json"
	"fmt"
)

// GetStashes Return a list of stashes path
func (s *Sensu) GetStashes() ([]interface{}, error) {
	return s.getList(fmt.Sprintf("stashes"), NoLimit)
}

// GetStash Get a stash
func (s *Sensu) GetStash(path string) (map[string]interface{}, error) {
	return s.getMap(fmt.Sprintf("stashes/%s", path))
}

// CreateStash create a stash (JSON document)
func (s *Sensu) CreateStash(payload interface{}) (map[string]interface{}, error) {
	payloadstr, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Stash parsing error: %q returned: %v", err, err)
	}
	return s.postPayload(fmt.Sprintf("stashes"), string(payloadstr[:]))
}

// CreateStashPath create a stash at path
func (s *Sensu) CreateStashPath(path string, payload map[string]interface{}) (map[string]interface{}, error) {
	payloadstr, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Stash parsing error: %q returned: %v", err, err)
	}
	return s.postPayload(fmt.Sprintf("stashes/%s", path), string(payloadstr[:]))
}

// DeleteStash Delete a stash (JSON document)
func (s *Sensu) DeleteStash(path string) error {
	return s.delete(fmt.Sprintf("stashes/%s", path))
}
