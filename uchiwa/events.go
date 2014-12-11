package uchiwa

import "github.com/palourde/logger"

// BuildEvents constructs events objects for frontend consumption
func BuildEvents() {
	for _, e := range tmpResults.Events {
		m := e.(map[string]interface{})

		// build backward compatible event object for Sensu < 0.13.0
		if m["id"] == nil {

			// build client object
			c := m["client"]
			delete(m, "client")
			m["client"] = map[string]interface{}{"name": c}

			// build check object
			c = m["check"]
			delete(m, "check")
			m["check"] = map[string]interface{}{"name": c, "issued": m["issued"], "output": m["output"], "status": m["status"], "occurrences": m["occurrences"]}

			// is flapping?
			if m["action"] == false {
				m["action"] = "create"
			} else {
				m["action"] = "flapping"
			}

			// remove old entries
			delete(m, "issued")
			delete(m, "output")
			delete(m, "status")
		}

		c, ok := m["client"].(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert event's client interface: %+v", c)
			continue
		}

		k := m["check"].(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert event's check interface: %+v", k)
			continue
		}

		m["acknowledged"] = isAcknowledged(c["name"].(string), k["name"].(string), m["dc"].(string))
	}
}

// ResolveEvent send a POST request to the /resolve endpoint in order to resolve an event
func ResolveEvent(data interface{}) error {

	api, m, err := findDcFromInterface(data)

	_, err = api.ResolveEvent(m["payload"])
	if err != nil {
		logger.Warning(err)
		return err
	}

	return nil
}
