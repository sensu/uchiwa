package uchiwa

import "github.com/sensu/uchiwa/uchiwa/logger"

// ResolveEvent send a POST request to the /resolve endpoint in order to resolve an event
func (u *Uchiwa) ResolveEvent(check, client, dc string) error {
	api, err := getAPI(u.Datacenters, dc)
	if err != nil {
		logger.Warning(err)
		return err
	}

	err = api.DeleteEvent(check, client)
	if err != nil {
		logger.Warning(err)
		return err
	}

	return nil
}
