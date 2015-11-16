package sensu

import "fmt"

// GetChecks returns a slice of all checks
func (s *Sensu) GetChecks() ([]interface{}, error) {
	return s.getSlice("checks", NoLimit)
}

// GetCheck returns a map of a specific check corresponding to the provided check name
func (s *Sensu) GetCheck(check string) (map[string]interface{}, error) {
	return s.getMap(fmt.Sprintf("checks/%s", check))
}
