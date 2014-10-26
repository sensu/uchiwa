package uchiwa

import (
	"errors"

	"github.com/palourde/logger"
)

// CreateStash send a POST request to the /stashes endpoint in order to create a stash
func CreateStash(data interface{}) error {

	api, m, err := findDcFromInterface(data)

	_, err = api.CreateStash(m["payload"])
	if err != nil {
		logger.Warning(err)
		return err
	}

	return nil
}

// DeleteStash send a DELETE request to the /stashes/*path* endpoint in order to delete a stash
func DeleteStash(data interface{}) error {
	api, m, err := findDcFromInterface(data)

	p, ok := m["payload"].(map[string]interface{})
	if !ok {
		logger.Warningf("Could not assert data interface %+v", data)
		return errors.New("Could not assert data interface")
	}

	err = api.DeleteStash(p["path"].(string))
	if err != nil {
		logger.Warning(err)
		return err
	}

	return nil
}
