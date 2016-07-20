package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

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

func (a *Config) login(user, pass string) (*User, error) {
	// Authenticate the user with the authentication driver
	u, err := a.DriverFn(user, pass)
	if err != nil {
		return nil, fmt.Errorf("Authentication failed: %s", err)
	}

	// Obfuscate the user's salt & hash
	u.PasswordHash = ""
	u.PasswordSalt = ""

	token, err := GetToken(&u.Role, user)
	if err != nil {
		return nil, fmt.Errorf("Authentication failed, could not create the token: %s", err)
	}

	// Add token to the user struct
	u.Token = token

	return u, nil
}
