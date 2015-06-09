package daemon

import (
	"errors"
	"fmt"

	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/uchiwa/sensu"
)

// FindDcFromInterface ...
func FindDcFromInterface(data interface{}, datacenters *[]sensu.Sensu) (*sensu.Sensu, map[string]interface{}, error) {
	m, ok := data.(map[string]interface{})
	if !ok {
		logger.Warningf("Type assertion failed. Could not assert the given interface into a map: %+v", data)
		return nil, nil, errors.New("Could not determine the datacenter.")
	}

	id := m["dc"].(string)
	if id == "" {
		logger.Warningf("The received interface does not contain any datacenter information: ", data)
		return nil, nil, errors.New("Could not determine the datacenter.")
	}

	for _, dc := range *datacenters {
		if dc.Name == id {
			return &dc, m, nil
		}
	}

	logger.Warningf("Could not find the datacenter %s into %+v: ", id, data)
	return nil, nil, fmt.Errorf("Could not find the datacenter %s", id)
}

// IsAcknowledged ...
func IsAcknowledged(client string, check string, dc string, stashes []interface{}) bool {
	if len(stashes) == 0 {
		return false
	}

	// add leading slash to check name
	if check != "" {
		check = fmt.Sprintf("/%s", check)
	}

	path := fmt.Sprintf("silence/%s%s", client, check)

	ack := false

	for _, stash := range stashes {
		m, ok := stash.(map[string]interface{})
		if !ok {
			continue
		}

		if m["path"] == path && m["dc"] == dc {
			ack = true
		}
	}

	return ack
}

func setDc(v interface{}, dc string) {
	m, ok := v.(map[string]interface{})
	if !ok {
		logger.Warningf("Could not assert interface: %+v", v)
	} else {
		m["dc"] = dc
	}
}

// StringInArray searches 'array' for 'item' string
// Returns true 'item' is a value of 'array'
func StringInArray(item string, array []string) bool {
	if item == "" || len(array) == 0 {
		return false
	}

	for _, element := range array {
		if element == item {
			return true
		}
	}

	return false
}
