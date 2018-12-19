package daemon

import (
	"errors"
	"fmt"

	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/sensu"
)

// FindDcFromInterface ...
func FindDcFromInterface(data interface{}, datacenters *[]sensu.Sensu) (*sensu.Sensu, map[string]interface{}, error) {
	m, ok := data.(map[string]interface{})
	if !ok {
		logger.Warningf("Type assertion failed. Could not assert the given interface into a map: %+v", data)
		return nil, nil, errors.New("could not determine the datacenter")
	}

	id := m["dc"].(string)
	if id == "" {
		logger.Warningf("The received interface does not contain any datacenter information: ", data)
		return nil, nil, errors.New("could not determine the datacenter")
	}

	for _, dc := range *datacenters {
		if dc.Name == id {
			return &dc, m, nil
		}
	}

	logger.Warningf("Could not find the datacenter %s into %+v: ", id, data)
	return nil, nil, fmt.Errorf("Could not find the datacenter %s", id)
}

// setID sets the _id attribute on every element of the slice from the dc and name
func setID(elements []interface{}, separator string) {
	for _, e := range elements {
		element, ok := e.(map[string]interface{})
		if !ok {
			continue
		}

		dc, ok := element["dc"].(string)
		if !ok {
			continue
		}

		name, ok := element["name"].(string)
		if !ok {
			// Support silence entries
			name, ok = element["id"].(string)
			if !ok {
				// Support stashes
				name, ok = element["path"].(string)
				if !ok {
					continue
				}
			}
		}

		element["_id"] = fmt.Sprintf("%s%s%s", dc, separator, name)
	}
}

func setDc(v interface{}, dc string) {
	m, ok := v.(map[string]interface{})
	if !ok {
		logger.Warningf("Could not assert interface: %+v", v)
	} else {
		m["dc"] = dc
	}
}
