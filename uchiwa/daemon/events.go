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
		clientMap, ok := m["client"].(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert event's client interface from %+v", clientMap)
			continue
		}

		client, ok := clientMap["name"].(string)
		if !ok {
			logger.Warningf("Could not assert event's client name from %+v", clientMap)
			continue
		}

		// get check name
		checkMap, ok := m["check"].(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert event's check interface from %+v", checkMap)
			continue
		}

		check, ok := checkMap["name"].(string)
		if !ok {
			logger.Warningf("Could not assert event's check name from %+v", checkMap)
			continue
		}

		// get dc name
		dc, ok := m["dc"].(string)
		if !ok {
			logger.Warningf("Could not assert event's datacenter name from %+v", m)
			continue
		}

		// Set the event unique ID
		m["_id"] = fmt.Sprintf("%s/%s/%s", dc, client, check)

		// Determine if the client is silenced
		m["client"].(map[string]interface{})["silenced"] = helpers.IsClientSilenced(client, dc, d.Data.Silenced)

		// Determine if the check is silenced.
		// See https://github.com/sensu/uchiwa/issues/602
		m["silenced"], m["silenced_by"] = helpers.IsCheckSilenced(checkMap, clientMap, dc, d.Data.Silenced)
	}
}
