package daemon

import (
	"fmt"

	"github.com/sensu/uchiwa/uchiwa/helpers"
	"github.com/sensu/uchiwa/uchiwa/logger"
)

// BuildEvents constructs events objects for frontend consumption
func (d *Daemon) buildEvents() {
	for _, e := range d.Data.Events {
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
			m["check"] = map[string]interface{}{"name": c, "issued": m["issued"], "output": m["output"], "status": m["status"]}

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

		// we assume the event isn't acknowledged in case we can't assert the following values
		m["acknowledged"] = false

		// get client name
		c, ok := m["client"].(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert event's client interface from %+v", c)
			continue
		}

		clientName, ok := c["name"].(string)
		if !ok {
			logger.Warningf("Could not assert event's client name from %+v", c)
			continue
		}

		// get check name
		k, ok := m["check"].(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert event's check interface from %+v", k)
			continue
		}

		checkName, ok := k["name"].(string)
		if !ok {
			logger.Warningf("Could not assert event's check name from %+v", k)
			continue
		}

		// get dc name
		dcName, ok := m["dc"].(string)
		if !ok {
			logger.Warningf("Could not assert event's datacenter name from %+v", m)
			continue
		}

		// Set the event unique ID
		m["_id"] = fmt.Sprintf("%s/%s/%s", dcName, clientName, checkName)

		// determine if the event is acknowledged
		m["acknowledged"] = helpers.IsAcknowledged(checkName, clientName, dcName, d.Data.Stashes)

		// detertermine if the client is acknowledged
		m["client"].(map[string]interface{})["acknowledged"] = helpers.IsAcknowledged("", clientName, dcName, d.Data.Stashes)
	}
}
