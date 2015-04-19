package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// Role ...
type Role struct {
	Get   bool
	Post  bool
	Admin bool
}

var adminRole = Role{Get: true, Post: true, Admin: true}
var operatorRole = Role{Get: true, Post: true, Admin: false}
var guestRole = Role{Get: true, Post: false, Admin: false}

func getRole(r string) (Role, error) {
	if r == "admin" {
		return adminRole, nil
	} else if r == "operator" {
		return operatorRole, nil
	} else if r == "guest" {
		return guestRole, nil
	}

	return Role{}, fmt.Errorf("No role '%s' found", r)
}

func hasPermission(t *jwt.Token, r *http.Request) bool {
	var role Role
	m := t.Claims["Role"]

	// use JSON representation of the interface to assert it into a Role struct
	j, _ := json.Marshal(&m)
	json.Unmarshal(j, &role)

	if r.URL.Path == "/users" && !role.Admin {
		return false
	} else if r.Method == "GET" {
		return true
	} else if role.Post || role.Admin {
		return true
	}

	return false
}
