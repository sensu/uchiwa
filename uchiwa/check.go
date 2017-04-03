package uchiwa

import (
	"fmt"

	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// GetCheck retrieves a specific check
func (u *Uchiwa) GetCheck(dc, name string) (map[string]interface{}, error) {
	api, err := getAPI(u.Datacenters, dc)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	check, err := api.GetCheck(name)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	// lock results
	//u.Mu.Lock()
	//defer u.Mu.Unlock()
	// client["_id"] = fmt.Sprintf("%s/%s", dc, name)
	// client["dc"] = dc
	// client["silenced"] = helpers.IsClientSilenced(name, dc, u.Data.Silenced)

	return check, nil
}

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
