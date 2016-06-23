package uchiwa

import (
	"fmt"

	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// IssueCheckExecution sends a POST request to the /stashes endpoint in order to create a stash
func (u *Uchiwa) IssueCheckExecution(data structs.CheckExecution) error {
	api, err := getAPI(u.Datacenters, data.Dc)
	if err != nil {
		logger.Warning(err)
		return err
	}

	_, err = api.IssueCheckExecution(data)
	if err != nil {
		logger.Warning(err)
		return err
	}

	return nil
}

func (u *Uchiwa) findCheck(name string) ([]interface{}, error) {
	var checks []interface{}
	for _, c := range u.Data.Checks {
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
