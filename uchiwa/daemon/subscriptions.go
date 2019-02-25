package daemon

import (
	"strings"

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
			logger.Debugf("%s", err)
			continue
		}

		for _, name := range generic.Subscriptions {
			// Do not add per-client subscriptions to the slice so we don't pollute
			// the subscriptions filter in the frontend.
			// See https://github.com/sensu/sensu-settings/pull/40.
			if strings.HasPrefix(strings.ToLower(name), "client:") {
				continue
			}

			subscription := structs.Subscription{Dc: generic.Dc, Name: name}

			if !isSubscriptionInSubscriptions(subscription, d.Data.Subscriptions) {
				d.Data.Subscriptions = append(d.Data.Subscriptions, subscription)
			}
		}
	}
}

func isSubscriptionInSubscriptions(subscription structs.Subscription, subscriptions []structs.Subscription) bool {
	for _, s := range subscriptions {
		if s.Dc == subscription.Dc && s.Name == subscription.Name {
			return true
		}
	}
	return false
}
