package sensu

import "fmt"

// DeleteCheckResult deletes a check result for a particular client
func (s *Sensu) DeleteCheckResult(check, client string) error {
	return s.delete(fmt.Sprintf("results/%s/%s", client, check))
}
