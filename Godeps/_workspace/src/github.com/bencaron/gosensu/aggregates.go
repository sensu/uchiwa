package sensu

import "fmt"

// GetAggregates Return the list of Aggregates
func (s *Sensu) GetAggregates() ([]interface{}, error) {
	return s.getList("aggregates", 0, 0)
}

// GetAggregatesSlice Return a slice in the current list of Aggregates
func (s *Sensu) GetAggregatesSlice(limit int, offset int) ([]interface{}, error) {
	return s.getList("aggregates", limit, offset)
}

// GetAggregate Return Aggregate info
func (s *Sensu) GetAggregate(check string, age int) ([]interface{}, error) {
	// FIXME GetAgregate Not handling age
	return s.getList(fmt.Sprintf("aggregate/%s", check), 0, 0)
}

// GetAggregateIssued Return Aggregate history
func (s *Sensu) GetAggregateIssued(check string, issued string, summarize bool, result bool) (map[string]interface{}, error) {
	// FIXME Aggregate Not handling summarize/result
	return s.get(fmt.Sprintf("aggregate/%s/%s", check, issued))
}

// DeleteAggregate Return the list of Aggregates
func (s *Sensu) DeleteAggregate(aggregate string) error {
	return s.delete(fmt.Sprintf("aggregate/%s", aggregate))
}
