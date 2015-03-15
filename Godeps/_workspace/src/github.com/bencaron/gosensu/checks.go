package sensu

import "fmt"

// GetChecks Return the list of checks
func (s *Sensu) GetChecks() ([]interface{}, error) {
	return s.getList("checks", 0, 0)
}

// GetCheck Return check info for a specific check
func (s *Sensu) GetCheck(check string) (map[string]interface{}, error) {
	return s.get(fmt.Sprintf("checks/%s", check))
}

// RequestCheck Issues a check request
func (s *Sensu) RequestCheck(checkName string) (map[string]interface{}, error) {
	return s.GetCheck(checkName)
	/*	rawcheck, err := s.GetCheck(checkName)
		if err != nil {
			return nil, fmt.Errorf("Can't RequestCheck for %s, error retrieving check: %s", checkName, err)
		}
		check := rawcheck.
		payload := fmt.Printf("{ \"check\": \"%s\", \"subscriber\": %v}", check["name"], json.Marshall(check["subscribers"]))
		return s.postPayload(fmt.Sprintf("check/request"))
	*/
}
