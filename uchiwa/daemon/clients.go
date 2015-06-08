package daemon

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// buildClients constructs clients objects for frontend consumption
func (d *Daemon) buildClients() {
	for _, c := range d.Data.Clients {
		client := c.(map[string]interface{})

		if client["version"] == nil {
			client["version"] = "0.12.x"
		}

		client = findClientEvents(client, &d.Data.Events)

		client["acknowledged"] = IsAcknowledged(client["name"].(string), "", client["dc"].(string), d.Data.Stashes)
	}
}

// findClientEvents searches for all events related to a particular client
// and set the status and output attributes of this client based on the events found
func findClientEvents(client map[string]interface{}, events *[]interface{}) map[string]interface{} {
	if len(*events) == 0 {
		client["status"] = 0
	} else {
		var criticals, warnings int
		var results []string
		for _, e := range *events {

			var event structs.GenericEvent
			err := mapstructure.Decode(e, &event)
			if err != nil {
				logger.Warningf("Could not convert the event to a generic event structure: %s", err)
				continue
			}

			// skip this event if not the right client
			if event.Client.Name != client["name"] || event.Dc != client["dc"] {
				continue
			}

			if event.Check.Status == 2 {
				criticals++
			} else if event.Check.Status == 1 {
				warnings++
			}

			results = append(results, event.Check.Output)
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

	return client
}
