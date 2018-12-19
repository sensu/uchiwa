package sensu

import (
	"context"
	"encoding/json"
	"fmt"
)

// ClearSilenced clears an entry from the silenced registry
func (s *Sensu) ClearSilenced(payload interface{}) (map[string]interface{}, error) {
	payloadstr, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Silence parsing error: %q returned: %v", err, err)
	}
	return s.postPayload(fmt.Sprintf("silenced/clear"), string(payloadstr[:]))
}

// GetSilenced returns the complete silenced registry
func (s *Sensu) GetSilenced(ctx context.Context) ([]interface{}, error) {
	return s.getSlice(ctx, fmt.Sprintf("silenced"), NoLimit)
}

// Silence updates the silenced registry with a new entry
func (s *Sensu) Silence(payload interface{}) (map[string]interface{}, error) {
	payloadstr, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Silence parsing error: %q returned: %v", err, err)
	}
	return s.postPayload(fmt.Sprintf("silenced"), string(payloadstr[:]))
}
