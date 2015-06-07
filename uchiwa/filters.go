package uchiwa

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/uchiwa/daemon"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// filterGetRequest is a function that filters a GET request that provides a 'get' attribute.
// The 'get' attribute is passed as a string argument
// Returns false is the request should not be filtered and therefore authorized
func filterGetRequest(dc string, token *jwt.Token) bool {
	if dc == "" {
		logger.Debug("The dc should not be empty")
		return true
	}

	if token == nil {
		logger.Debug("No token found in the request, no filters will be applied")
		return false
	}

	role, err := getRoleFromToken(token)
	if err != nil {
		logger.Warningf("%s", err)
		return true
	}

	if len(role.Datacenters) == 0 || daemon.StringInArray(dc, role.Datacenters) {
		return false
	}

	return true
}

// filterPostRequest is a function that filters a POST request.
// The JWT and the data posted are passed as arguments
// Returns false if the request should not be filtered and therefore authorized
func filterPostRequest(token *jwt.Token, data *interface{}) bool {
	if token == nil {
		logger.Debug("No token found in the request, no filters will be applied")
		return false
	}

	role, err := getRoleFromToken(token)
	if err != nil {
		logger.Warningf("%s", err)
		return true
	}

	// do not filter if both datacenters & subscriptions filters are empty
	if len(role.Datacenters) == 0 && len(role.Subscriptions) == 0 {
		logger.Debugf("No datacenter and subscription filters found in the role %s", role.Name)
		return false
	}

	// decode the data interface to a generic event structure
	var generic structs.GenericEvent
	err = mapstructure.Decode(*data, &generic)
	if err != nil {
		logger.Debug("%s", err)
		return true
	}

	if len(role.Datacenters) == 0 || daemon.StringInArray(generic.Dc, role.Datacenters) {
		if len(role.Subscriptions) == 0 || len(generic.Check.Subscribers) == 0 {
			return false
		} else if sliceIntersection(generic.Check.Subscribers, role.Subscriptions) {
			return false
		}
	}

	return true
}

// filterSensu is a function that filters Sensu Data based on
// the datacenters and subscriptions within the Role struct of the JWT
func filterSensu(token *jwt.Token, data *structs.Data) *structs.Data {
	if token == nil {
		logger.Debug("No token found in the request, returning all data")
		return data
	}

	role, err := getRoleFromToken(token)
	if err != nil {
		logger.Warningf("%s", err)
		return &structs.Data{}
	}

	// return all data if no datacenters are found
	if len(role.Datacenters) == 0 && len(role.Subscriptions) == 0 {
		logger.Debugf("No datacenter and subscription filters found in the role %s", role.Name)
		return data
	}

	var filteredData structs.Data

	// Aggregates
	for _, aggregate := range data.Aggregates {
		var generic structs.Generic
		err := mapstructure.Decode(aggregate, &generic)
		if err != nil {
			continue
		}

		// verify if the generic element is part of the datacenters specified within the role
		if len(role.Datacenters) == 0 || daemon.StringInArray(generic.Dc, role.Datacenters) {
			filteredData.Aggregates = append(filteredData.Aggregates, aggregate)
		}
	}

	// Checks
	for _, check := range data.Checks {
		var generic structs.GenericCheck
		err := mapstructure.Decode(check, &generic)
		if err != nil {
			continue
		}

		// verify if the generic element is part of the datacenters and the subscriptions specified within the role
		if len(role.Datacenters) == 0 || daemon.StringInArray(generic.Dc, role.Datacenters) {
			if len(role.Subscriptions) == 0 || sliceIntersection(generic.Subscribers, role.Subscriptions) {
				filteredData.Checks = append(filteredData.Checks, check)
			}
		}

	}

	// Clients
	for _, client := range data.Clients {
		var generic structs.GenericClient
		err := mapstructure.Decode(client, &generic)
		if err != nil {
			continue
		}

		// verify if the generic element is part of the datacenters and the subscriptions specified within the role
		if len(role.Datacenters) == 0 || daemon.StringInArray(generic.Dc, role.Datacenters) {
			if len(role.Subscriptions) == 0 || sliceIntersection(generic.Subscriptions, role.Subscriptions) {
				filteredData.Clients = append(filteredData.Clients, client)
			}
		}
	}

	// Events
	for _, event := range data.Events {
		var generic structs.GenericEvent
		err := mapstructure.Decode(event, &generic)
		if err != nil {
			continue
		}

		// verify if the generic element is part of the datacenters and the subscriptions specified within the role
		if len(role.Datacenters) == 0 || daemon.StringInArray(generic.Dc, role.Datacenters) {
			if len(role.Subscriptions) == 0 || sliceIntersection(generic.Check.Subscribers, role.Subscriptions) {
				filteredData.Events = append(filteredData.Events, event)
			}
		}
	}

	// Stashes
	for _, stash := range data.Stashes {
		var generic structs.Generic
		err := mapstructure.Decode(stash, &generic)
		if err != nil {
			continue
		}

		// verify if the generic element is part of the datacenters specified within the role
		if len(role.Datacenters) == 0 || daemon.StringInArray(generic.Dc, role.Datacenters) {
			filteredData.Stashes = append(filteredData.Stashes, stash)
		}
	}

	// Datacenters
	for _, datacenter := range data.Dc {
		// verify if the datacenter is part of the datacenters specified within the role
		if len(role.Datacenters) == 0 || daemon.StringInArray(datacenter.Name, role.Datacenters) {
			filteredData.Dc = append(filteredData.Dc, datacenter)
		}
	}

	// Subscriptions
	for _, subscription := range data.Subscriptions {
		if len(role.Subscriptions) == 0 || daemon.StringInArray(subscription, role.Subscriptions) {
			filteredData.Subscriptions = append(filteredData.Subscriptions, subscription)
		}
	}

	return &filteredData
}
