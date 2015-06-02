package sensu

import "fmt"

// GetChecks Return the list of checks
func (s *Sensu) GetChecks() ([]interface{}, error) {
	return s.getList("checks", 0, 0)
}

// GetCheck Return check info for a specific check
func (s *Sensu) GetCheck(check string) (map[string]interface{}, error) {
	return s.getMap(fmt.Sprintf("checks/%s", check))
}
