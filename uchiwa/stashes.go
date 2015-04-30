package uchiwa

import "github.com/palourde/logger"

type stash struct {
	Dc      string      `json:"dc"`
	Path    string      `json:"path"`
	Content interface{} `json:"content"`
	Expire  int32       `json:"expire"`
}

// PostStash send a POST request to the /stashes endpoint in order to create a stash
func PostStash(data stash) error {
	api, err := getAPI(data.Dc)
	if err != nil {
		logger.Warning(err)
		return err
	}

	_, err = api.CreateStash(data)
	if err != nil {
		logger.Warning(err)
		return err
	}

	return nil
}

// DeleteStash send a DELETE request to the /stashes/*path* endpoint in order to delete a stash
func DeleteStash(data stash) error {
	api, err := getAPI(data.Dc)
	if err != nil {
		logger.Warning(err)
		return err
	}

	err = api.DeleteStash(data.Path)
	if err != nil {
		logger.Warning(err)
		return err
	}

	return nil
}
