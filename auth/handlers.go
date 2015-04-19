package auth

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/palourde/logger"
)

// publicHandler enforce no authentication
func publicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// restrictedHandler enforce authentication by validating a token
func restrictedHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := jwt.ParseFromRequest(r, func(t *jwt.Token) (interface{}, error) {
			return pubKeyPEM, nil
		})
		if token != nil && err == nil {
			if token.Valid {
				authorized := hasPermission(token, r)
				if !authorized {
					http.Error(w, "", http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r)
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
		user, err := a.Driver(u, p)
		if err != nil {
			logger.Infof("Authentication failed: %s", err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// obfuscate the user's salt & hash
		user.PasswordHash = ""
		user.PasswordSalt = ""

		role, err := getRole(user.Role)
		if err != nil {
			logger.Infof("%s", err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		token, err := GetToken(role)
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
