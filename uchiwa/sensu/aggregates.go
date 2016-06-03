package sensu

import "fmt"

// GetAggregates returns a slice of all aggregates
func (s *Sensu) GetAggregates() ([]interface{}, error) {
	return s.getSlice("aggregates", NoLimit)
}

// GetAggregate returns a map of a specific aggregate corresponding to the provided check name
// Only supports Sensu >= 0.24.0
func (s *Sensu) GetAggregate(check string) (map[string]interface{}, error) {
	return s.getMap(fmt.Sprintf("aggregates/%s", check))
}

// GetAggregateIssued returns a map containing the history of a specific check corresponding to the provided check name and the issued timestamp
func (s *Sensu) GetAggregateIssued(check string, issued string) (map[string]interface{}, error) {
	return s.getMap(fmt.Sprintf("aggregates/%s/%s", check, issued))
}

// DeleteAggregate deletes an aggregate using its check name
func (s *Sensu) DeleteAggregate(check string) error {
	return s.delete(fmt.Sprintf("aggregates/%s", check))
}
