package uchiwa

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/uchiwa/auth"
	"github.com/sensu/uchiwa/uchiwa/sensu"
)

func getAPI(datacenters *[]sensu.Sensu, name string) (*sensu.Sensu, error) {
	if name == "" {
		return nil, errors.New("The datacenter name can't be empty")
	}

	for _, datacenter := range *datacenters {
		if datacenter.Name == name {
			return &datacenter, nil
		}
	}

	return nil, fmt.Errorf("Could not find the datacenter '%s'", name)
}

func findModel(id string, dc string, checks []interface{}) map[string]interface{} {
	for _, k := range checks {
		m, ok := k.(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert check interface %+v", k)
			continue
		}
		if m["name"] == id && m["dc"] == dc {
			return m
		}
	}
	return nil
}

func getRoleFromToken(token *jwt.Token) (*auth.Role, error) {
	r, ok := token.Claims["Role"]
	if !ok {
		return &auth.Role{}, errors.New("Could not retrieve the user Role from the JWT")
	}

	var role auth.Role
	err := mapstructure.Decode(r, &role)
	if err != nil {
		return &auth.Role{}, err
	}

	return &role, nil
}

// arrayIntersection searches for values in both arrays
// Returns true if there's at least one intersection
func arrayIntersection(array1, array2 []string) bool {
	if len(array1) == 0 || len(array2) == 0 {
		return false
	}

	for _, a := range array1 {
		for _, b := range array2 {
			if a == b {
				return true
			}
		}
	}

	return false
}
