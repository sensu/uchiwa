package uchiwa

import (
	"errors"
	"fmt"

	"github.com/bencaron/gosensu"
	"github.com/palourde/logger"
)

func findDcFromInterface(data interface{}) (*sensu.Sensu, map[string]interface{}, error) {
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

	for _, dc := range datacenters {
		if dc.Name == id {
			return &dc, m, nil
		}
	}

	logger.Warningf("Could not find the datacenter %s into %+v: ", id, data)
	return nil, nil, fmt.Errorf("Could not find the datacenter %s", id)
}

func findDcFromString(id *string) (*sensu.Sensu, error) {
	for _, d := range datacenters {
		if d.Name == *id {
			return &d, nil
		}
	}
	return nil, fmt.Errorf("Could not find datacenter %s", *id)
}

func findModel(id string, dc string) map[string]interface{} {
	for _, k := range tmpResults.Checks {
		m, ok := k.(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert check interface %+v", k)
			continue
		}
		if m["name"] == id {
			return m
		}
	}
	return nil
}

func findStatus(client map[string]interface{}) {
	if len(tmpResults.Events) == 0 {
		client["status"] = 0
	} else {
		var criticals, warnings int
		var results []string
		for _, e := range tmpResults.Events {
			m, ok := e.(map[string]interface{})
			if !ok {
				logger.Warningf("Could not assert event interface %+v", e)
				continue
			}

			// skip this event if another dc
			if m["dc"] != client["dc"] {
				continue
			}

			c, ok := m["client"].(map[string]interface{})
			if !ok {
				logger.Warningf("Could not assert event's client interface: %+v", c)
				continue
			}

			// skip this event if another client
			if c["name"] != client["name"] || m["dc"] != client["dc"] {
				continue
			}

			k := m["check"].(map[string]interface{})
			if !ok {
				logger.Warningf("Could not assert event's check interface: %+v", k)
				continue
			}

			results = append(results, k["output"].(string))

			status := int(k["status"].(float64))
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

func isAcknowledged(c string, k string, dc string) bool {
	if len(tmpResults.Stashes) == 0 {
		return false
	}

	// add leading slash to check name
	if k != "" {
		k = fmt.Sprintf("/%s", k)
	}

	p := fmt.Sprintf("silence/%s%s", c, k)

	a := false

	for _, s := range tmpResults.Stashes {
		m, ok := s.(map[string]interface{})
		if !ok {
			continue
		}

		if m["path"] == p && m["dc"] == dc {
			a = true
		}
	}

	return a
}

func setDc(v interface{}, dc string) {
	m, ok := v.(map[string]interface{})
	if !ok {
		logger.Warningf("Could not assert interface: %+v", v)
	} else {
		m["dc"] = dc
	}
}
