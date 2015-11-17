package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/sensu/uchiwa/uchiwa/audit"
	"github.com/sensu/uchiwa/uchiwa/helpers"
	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

const jwtToken = "jwtToken"

// GetTokenFromContext retrieves the JWT Token from the request
func GetTokenFromContext(r *http.Request) *jwt.Token {
	if value := context.Get(r, jwtToken); value != nil {
		return value.(*jwt.Token)
	}
	return nil
}

// setTokenIntoContext inject the JWT Token into the request for later use
func setTokenIntoContext(r *http.Request, token *jwt.Token) {
	context.Set(r, jwtToken, token)
}

// publicHandler does not enforce authentication
func publicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// restrictedHandler enforce authentication by validating the JWT token
func restrictedHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseFromRequest(r, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}

			return publicKey, nil
		})
		if token != nil && err == nil {
			if token.Valid {
				setTokenIntoContext(r, token)
				authorized := hasPermission(token, r)
				if !authorized {
					http.Error(w, "", http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r)
				context.Clear(r)
			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	})
}

// Authenticate calls the proper handler based on whether authentication is enabled or not
func (a *Config) Authenticate(next http.Handler) http.Handler {
	if a.DriverName == "none" {
		return publicHandler(next)
	}
	return restrictedHandler(next)
}

// GetIdentification retrieves the user & pass from a POST and authenticates the user against the Identification driver
func (a *Config) GetIdentification() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Redirect(w, r, "/#/login", http.StatusFound)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var data interface{}
		err := decoder.Decode(&data)
		if err != nil {
			logger.Warningf("Could not decode the body: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		m, ok := data.(map[string]interface{})
		if !ok {
			logger.Warningf("Could not assert the body: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		u := m["user"].(string)
		p := m["pass"].(string)
		if u == "" || p == "" {
			logger.Info("Authentication failed: user and password must not be empty")
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// validate the user with the Login authentication driver
		user, err := a.DriverFn(u, p)
		if err != nil {
			message := fmt.Sprintf("Authentication failed: %s", err)

			// Output to stdout
			logger.Info(message)

			// Output to audit log
			log := structs.AuditLog{Action: "loginfailure", Level: "default", Output: message}
			log.RemoteAddr = helpers.GetIP(r)
			audit.Log(log)

			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// obfuscate the user's salt & hash
		user.PasswordHash = ""
		user.PasswordSalt = ""

		token, err := GetToken(&user.Role, u)
		if err != nil {
			logger.Warningf("Authentication failed, could not create the token: %s", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Add token to the user struct
		user.Token = token

		j, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		return
	})
}
