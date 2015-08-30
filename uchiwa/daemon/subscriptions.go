package daemon

import (
	"github.com/mitchellh/mapstructure"
	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// BuildSubscriptions builds a slice of every client subscriptions
func (d *Daemon) BuildSubscriptions() {
	for _, client := range d.Data.Clients {
		var generic structs.GenericClient
		err := mapstructure.Decode(client, &generic)
		if err != nil {
			logger.Debug("%s", err)
			continue
		}

		for _, subscription := range generic.Subscriptions {
			if !StringInArray(subscription, d.Data.Subscriptions) {
				d.Data.Subscriptions = append(d.Data.Subscriptions, subscription)
			}
		}
	}
}
