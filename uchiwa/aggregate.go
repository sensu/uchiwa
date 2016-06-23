package uchiwa

import (
	"fmt"

	"github.com/sensu/uchiwa/uchiwa/logger"
)

// GetAggregate retrieves a list of issued timestamps from a specified DC
func (u *Uchiwa) GetAggregate(check string, dc string) (*map[string]interface{}, error) {
	api, err := getAPI(u.Datacenters, dc)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	aggregate, err := api.GetAggregate(check)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	return &aggregate, nil
}

// GetAggregateByIssued retrieves aggregate check info from a specified DC
func (u *Uchiwa) GetAggregateByIssued(check string, issued string, dc string) (*map[string]interface{}, error) {
	api, err := getAPI(u.Datacenters, dc)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	aggregate, err := api.GetAggregateIssued(check, issued)
	if err != nil {
		return nil, err
	}

	return &aggregate, nil
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
