package uchiwa

import (
	"errors"
	"fmt"

	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/uchiwa/sensu"
)

func getAPI(datacenters *[]sensu.Sensu, name string) (*sensu.Sensu, error) {
	if name == "" {
		return nil, errors.New("The datacenter name can't be empty")
	}

	for _, datacenter := range *datacenters {
		if datacenter.Name == name {
			return &datacenter, nil
		}
	}

	return nil, fmt.Errorf("Could not find the datacenter '%s'", name)
}

func findModel(id string, dc string, checks []interface{}) map[string]interface{} {
	for _, k := range checks {
		m, ok := k.(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert check interface %+v", k)
			continue
		}
		if m["name"] == id && m["dc"] == dc {
			return m
		}
	}
	return nil
}

// inArray searches 'array' for 'item'
// Returns true if 'array' is empty
func inArray(item string, array []string) bool {
	if len(array) == 0 {
		return true
	}

	if item == "" {
		return false
	}

	for _, element := range array {
		if element == item {
			return true
		}
	}

	return false
}
