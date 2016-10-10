package uchiwa

import (
	"fmt"

	"github.com/sensu/uchiwa/uchiwa/logger"
)

// GetAggregate retrieves a specific aggregate
func (u *Uchiwa) GetAggregate(name, dc string) (*map[string]interface{}, error) {
	api, err := getAPI(u.Datacenters, dc)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	aggregate, err := api.GetAggregate(name)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	return &aggregate, nil
}

// GetAggregateChecks retrieves check members of an aggregate
func (u *Uchiwa) GetAggregateChecks(name, dc string) (*[]interface{}, error) {
	api, err := getAPI(u.Datacenters, dc)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	checks, err := api.GetAggregateChecks(name)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	return &checks, nil
}

// GetAggregateClients retrieves client members of an aggregate
func (u *Uchiwa) GetAggregateClients(name, dc string) (*[]interface{}, error) {
	api, err := getAPI(u.Datacenters, dc)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	clients, err := api.GetAggregateClients(name)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	return &clients, nil
}

// GetAggregateResults retrieves check result members by severity of an aggregate
func (u *Uchiwa) GetAggregateResults(name, severity, dc string) (*[]interface{}, error) {
	api, err := getAPI(u.Datacenters, dc)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	results, err := api.GetAggregateResults(name, severity)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	return &results, nil
}

func (u *Uchiwa) findAggregate(name string) ([]interface{}, error) {
	var checks []interface{}
	for _, c := range u.Data.Aggregates {
		m, ok := c.(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert this check to an interface %+v", c)
			continue
		}
		if m["name"] == name {
			checks = append(checks, m)
		}
	}

	if len(checks) == 0 {
		return nil, fmt.Errorf("Could not find any checks with the name '%s'", name)
	}

	return checks, nil
}
