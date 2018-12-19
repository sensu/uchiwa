package sensu

import (
	"context"
	"encoding/json"
	"fmt"
)

// GetChecks returns a slice of all checks
func (s *Sensu) GetChecks(ctx context.Context) ([]interface{}, error) {
	return s.getSlice(ctx, "checks", NoLimit)
}

// GetCheck returns a map of a specific check corresponding to the provided check name
func (s *Sensu) GetCheck(check string) (map[string]interface{}, error) {
	return s.getMap(fmt.Sprintf("checks/%s", check))
}

// IssueCheckExecution send a POST request to the /request endpoint in order
// to issue a check execution request
func (s *Sensu) IssueCheckExecution(payload interface{}) (map[string]interface{}, error) {
	payloadstr, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("Stash parsing error: %q returned: %v", err, err)
	}
	return s.postPayload(fmt.Sprintf("request"), string(payloadstr[:]))
}
