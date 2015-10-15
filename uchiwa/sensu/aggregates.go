package sensu

import "fmt"

// GetAggregates Return the list of Aggregates
func (s *Sensu) GetAggregates() ([]interface{}, error) {
	return s.getList("aggregates", NoLimit)
}

// GetAggregate Return Aggregate info
func (s *Sensu) GetAggregate(check string, age int) ([]interface{}, error) {
	// TODO GetAgregate Not handling age
	return s.getList(fmt.Sprintf("aggregate/%s", check), NoLimit)
}

// GetAggregateIssued Return Aggregate history
func (s *Sensu) GetAggregateIssued(check string, issued string, summarize bool, result bool) (map[string]interface{}, error) {
	// TODO Aggregate Not handling summarize/result
	return s.getMap(fmt.Sprintf("aggregate/%s/%s", check, issued))
}

// DeleteAggregate Return the list of Aggregates
func (s *Sensu) DeleteAggregate(aggregate string) error {
	return s.delete(fmt.Sprintf("aggregate/%s", aggregate))
}
