package uchiwa

import "github.com/sensu/uchiwa/uchiwa/logger"

// DeleteCheckResult sends a DELETE request in order to
// remove the result for a given check on a given client
func (u *Uchiwa) DeleteCheckResult(check, client, dc string) error {
	api, err := getAPI(u.Datacenters, dc)
	if err != nil {
		logger.Warning(err)
		return err
	}

	err = api.DeleteCheckResult(check, client)
	if err != nil {
		logger.Warning(err)
		return err
	}

	return nil
}
