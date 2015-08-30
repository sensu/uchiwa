package uchiwa

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/sensu/uchiwa/uchiwa/logger"
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

// GetIP returns the real user IP address
func GetIP(r *http.Request) string {
	if xForwardedFor := r.Header.Get("X-FORWARDED-FOR"); len(xForwardedFor) > 0 {
		return xForwardedFor
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// GetRoleFromToken ...
func GetRoleFromToken(token *jwt.Token) (*auth.Role, error) {
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

// MergeStringSlices merges two slices of strings and remove duplicated values
func MergeStringSlices(a1, a2 []string) []string {
	if len(a1) == 0 {
		return a2
	} else if len(a2) == 0 {
		return a1
	}

	s := make([]string, len(a1), len(a1)+len(a2))
	copy(s, a1)

next:
	for _, x := range a2 {
		for _, y := range s {
			if x == y {
				continue next
			}
		}
		s = append(s, x)
	}
	return s
}

// SliceIntersection searches for values in both slices
// Returns true if there's at least one intersection
func SliceIntersection(a1, a2 []string) bool {
	if len(a1) == 0 || len(a2) == 0 {
		return false
	}

	for _, x := range a1 {
		for _, y := range a2 {
			if x == y {
				return true
			}
		}
	}

	return false
}
