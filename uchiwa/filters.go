package uchiwa

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/uchiwa/auth"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// filterSensu
func filterSensu(token *jwt.Token, data *structs.Data) *structs.Data {
	if token == nil {
		logger.Debug("No token found in the request, returning all data")
		return data
	}

	r, ok := token.Claims["Role"]
	if !ok {
		logger.Warning("Could not retrieve the user role from the token")
		return &structs.Data{}
	}

	var role auth.Role
	err := mapstructure.Decode(r, &role)
	if err != nil {
		logger.Warning(err)
		return &structs.Data{}
	}

	filteredData := findDatacenter(&role, data)

	//fmt.Println(filteredData)
	return filteredData
}

func findDatacenter(role *auth.Role, data *structs.Data) *structs.Data {
	// return all data if no datacenters are found
	if len(role.Datacenters) == 0 {
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

		// check if the generic element is part of the datacenters specified within the role
		if isMemberOfDatacenter(role.Datacenters, generic.Dc) {
			filteredData.Aggregates = append(filteredData.Aggregates, aggregate)
		}
	}

	// Checks
	for _, check := range data.Checks {
		var generic structs.Generic
		err := mapstructure.Decode(check, &generic)
		if err != nil {
			continue
		}

		// check if the generic element is part of the datacenters specified within the role
		if isMemberOfDatacenter(role.Datacenters, generic.Dc) {
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

		// check if the generic element is part of the datacenters specified within the role
		if isMemberOfDatacenter(role.Datacenters, generic.Dc) {
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

		// check if the generic element is part of the datacenters specified within the role
		if isMemberOfDatacenter(role.Datacenters, generic.Dc) {
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

		// check if the generic element is part of the datacenters specified within the role
		if isMemberOfDatacenter(role.Datacenters, generic.Dc) {
			filteredData.Stashes = append(filteredData.Stashes, stash)
		}
	}

	return &filteredData
}

func isMemberOfDatacenter(datacenters []string, name string) bool {
	if name == "" {
		return false
	}

	for _, datacenter := range datacenters {
		if datacenter == name {
			return true
		}
	}

	return false
}
