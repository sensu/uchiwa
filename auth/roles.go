package auth

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// Role contains the roles of each GitHub team
type Role struct {
	Name     string
	Members  []string
	Readonly bool
}

func hasPermission(t *jwt.Token, r *http.Request) bool {
	var role Role
	m := t.Claims["Role"]

	// use JSON representation of the interface to assert it into the uchiwa.Role struct
	j, _ := json.Marshal(&m)
	json.Unmarshal(j, &role)

	if r.Method == "GET" {
		return true
	} else if !role.Readonly {
		return true
	}

	return false
}
