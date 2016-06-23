package uchiwa

import (
	"fmt"

	"github.com/sensu/uchiwa/uchiwa/logger"
)

type stash struct {
	Dc      string                 `json:"dc"`
	Path    string                 `json:"path"`
	Content map[string]interface{} `json:"content"`
	Expire  int32                  `json:"expire,omitempty"`
}

// PostStash send a POST request to the /stashes endpoint in order to create a stash
func (u *Uchiwa) PostStash(data stash) error {
	api, err := getAPI(u.Datacenters, data.Dc)
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
func (u *Uchiwa) DeleteStash(dc, path string) error {
	api, err := getAPI(u.Datacenters, dc)
	if err != nil {
		return err
	}

	err = api.DeleteStash(path)
	if err != nil {
		return err
	}

	return nil
}

func (u *Uchiwa) findStash(path string) ([]interface{}, error) {
	var stashes []interface{}
	for _, c := range u.Data.Stashes {
		m, ok := c.(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert this stash to an interface %+v", c)
			continue
		}
		if m["path"] == path {
			stashes = append(stashes, m)
		}
	}

	if len(stashes) == 0 {
		return nil, fmt.Errorf("Could not find any stashes with the path '%s'", path)
	}

	return stashes, nil
}
