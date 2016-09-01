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

		// detertermine if the client is acknowledged
		m["client"].(map[string]interface{})["silenced"] = helpers.IsClientSilenced(clientName, dcName, d.Data.Silenced)
	}
}
