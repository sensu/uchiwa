package sensu

import (
	"context"
	"fmt"
)

// DeleteAggregate deletes an aggregate using its check name
func (s *Sensu) DeleteAggregate(name string) error {
	return s.delete(fmt.Sprintf("aggregates/%s", name))
}

// GetAggregates returns a slice of all aggregates
func (s *Sensu) GetAggregates(ctx context.Context) ([]interface{}, error) {
	return s.getSlice(ctx, "aggregates", NoLimit)
}

// GetAggregate returns a map of a specific aggregate corresponding to the provided check name
func (s *Sensu) GetAggregate(name string) (map[string]interface{}, error) {
	return s.getMap(fmt.Sprintf("aggregates/%s", name))
}

// GetAggregateChecks returns a slice of all checks members of an aggregate
func (s *Sensu) GetAggregateChecks(name string) ([]interface{}, error) {
	return s.getSlice(context.Background(), fmt.Sprintf("aggregates/%s/checks", name), NoLimit)
}

// GetAggregateClients returns a slice of all clients members of an aggregate
func (s *Sensu) GetAggregateClients(name string) ([]interface{}, error) {
	return s.getSlice(context.Background(), fmt.Sprintf("aggregates/%s/clients", name), NoLimit)
}

// GetAggregateResults returns a slice of all check result members by severity
func (s *Sensu) GetAggregateResults(name, severity string) ([]interface{}, error) {
	return s.getSlice(context.Background(), fmt.Sprintf("aggregates/%s/results/%s", name, severity), NoLimit)
}
