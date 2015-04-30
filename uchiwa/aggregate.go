package uchiwa

import (
	"github.com/palourde/logger"
)

// GetAggregate retrieves a list of issued timestamps from a specified DC
func GetAggregate(check string, dc string) (*[]interface{}, error) {
	api, err := getAPI(dc)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	aggregate, err := api.GetAggregate(check, 1)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	return &aggregate, nil
}

// GetAggregateByIssued retrieves aggregate check info from a specified DC
func GetAggregateByIssued(check string, issued string, dc string) (*map[string]interface{}, error) {
	api, err := getAPI(dc)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	aggregate, err := api.GetAggregateIssued(check, issued, true, true)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	return &aggregate, nil
}
