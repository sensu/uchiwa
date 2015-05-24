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

func (d *Daemon) findStatus(client map[string]interface{}) {
	if len(d.Data.Events) == 0 {
		client["status"] = 0
	} else {
		var criticals, warnings int
		var results []string
		for _, event := range d.Data.Events {
			m, ok := event.(map[string]interface{})
			if !ok {
				logger.Warningf("Could not assert the event %+v", event)
				continue
			}

			// skip this event if another dc
			if m["dc"] != client["dc"] {
				continue
			}

			c, ok := m["client"].(map[string]interface{})
			if !ok {
				logger.Warningf("Could not assert event's client: %+v", c)
				continue
			}

			// skip this event if another client
			if c["name"] != client["name"] || m["dc"] != client["dc"] {
				continue
			}

			check := m["check"].(map[string]interface{})
			if !ok {
				logger.Warningf("Could not assert event's check interface: %+v", check)
				continue
			}

			results = append(results, check["output"].(string))

			status := int(check["status"].(float64))
			if status == 2 {
				criticals++
			} else if status == 1 {
				warnings++
			}
		}

		if len(results) == 0 {
			client["status"] = 0
		} else if criticals > 0 {
			client["status"] = 2
		} else if warnings > 0 {
			client["status"] = 1
		} else {
			client["status"] = 3
		}

		if len(results) == 1 {
			client["output"] = results[0]
		} else if len(results) > 1 {
			output := fmt.Sprintf("%s and %d more...", results[0], (len(results) - 1))
			client["output"] = output
		}
	}
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
