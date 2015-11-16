package sensu

import "fmt"

// GetAggregates returns a slice of all aggregates
func (s *Sensu) GetAggregates() ([]interface{}, error) {
	return s.getSlice("aggregates", NoLimit)
}

// GetAggregate returns a slice of a specific aggregate corresponding to the provided check name
func (s *Sensu) GetAggregate(check string, age int) ([]interface{}, error) {
	return s.getSlice(fmt.Sprintf("aggregate/%s", check), NoLimit)
}

// GetAggregateIssued returns a map containing the history of a specific check corresponding to the provided check name and the issued timestamp
func (s *Sensu) GetAggregateIssued(check string, issued string) (map[string]interface{}, error) {
	return s.getMap(fmt.Sprintf("aggregate/%s/%s", check, issued))
}

// DeleteAggregate deletes an aggregate using its check name
func (s *Sensu) DeleteAggregate(check string) error {
	return s.delete(fmt.Sprintf("aggregate/%s", check))
}
