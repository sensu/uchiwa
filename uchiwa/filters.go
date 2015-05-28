package uchiwa

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/uchiwa/auth"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// filterGetSensu is a function that filters Sensu Data based on
// the datacenters and subscriptions within the Role struct of the JWT
func filterGetSensu(token *jwt.Token, data *structs.Data) *structs.Data {
	if token == nil {
		logger.Debug("No token found in the request, returning all data")
		return data
	}

	r, ok := token.Claims["Role"]
	if !ok {
		logger.Warning("Could not retrieve the user Role from the JWT")
		return &structs.Data{}
	}

	var role auth.Role
	err := mapstructure.Decode(r, &role)
	if err != nil {
		logger.Warning(err)
		return &structs.Data{}
	}

	// return all data if no datacenters are found
	if len(role.Datacenters) == 0 && len(role.Subscriptions) == 0 {
		logger.Debugf("No datacenters found in the role %s", role.Name)
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
		if inArray(generic.Dc, role.Datacenters) {
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
		fmt.Println(generic)
		// verify if the generic element is part of the datacenters specified within the role
		if inArray(generic.Dc, role.Datacenters) {
			filteredData.Checks = append(filteredData.Checks, check)
		}
	}

	// Clients
	for _, client := range data.Clients {
		var generic structs.Generic
		err := mapstructure.Decode(client, &generic)
		if err != nil {
			continue
		}

		// verify if the generic element is part of the datacenters specified within the role
		if inArray(generic.Dc, role.Datacenters) {
			filteredData.Clients = append(filteredData.Clients, client)
		}
	}

	// Events
	for _, event := range data.Events {
		var generic structs.Generic
		err := mapstructure.Decode(event, &generic)
		if err != nil {
			continue
		}

		// verify if the generic element is part of the datacenters specified within the role
		if inArray(generic.Dc, role.Datacenters) {
			filteredData.Events = append(filteredData.Events, event)
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
		if inArray(generic.Dc, role.Datacenters) {
			filteredData.Stashes = append(filteredData.Stashes, stash)
		}
	}

	// Datacenters
	for _, datacenter := range data.Dc {
		// verify if the datacenter is part of the datacenters specified within the role
		if inArray(datacenter.Name, role.Datacenters) {
			filteredData.Dc = append(filteredData.Dc, datacenter)
		}
	}

	//fmt.Println(filteredData)
	return &filteredData
}
