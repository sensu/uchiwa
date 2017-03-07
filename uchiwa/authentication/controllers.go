package authentication

import (
	"encoding/json"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/sensu/uchiwa/uchiwa/audit"
	"github.com/sensu/uchiwa/uchiwa/helpers"
	"github.com/sensu/uchiwa/uchiwa/logger"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// New function initalizes and returns a Config struct
func New(auth structs.Auth) Config {
	c := Config{
		Auth: auth,
	}
	return c
}

// publicHandler does not enforce authentication
func publicHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// restrictedHandler enforce authentication by validating the JWT
// or the access token provided in the configuration
func restrictedHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token *jwt.Token
		authenticationToken, err := r.Cookie(authenticationCookieName)
		if err == nil {
			// Verify the JWT
			token, err = verifyJWT(authenticationToken.Value)
			if err != nil {
				logger.Debug("No JWT token provided")
			} else {
				xsrfTokenFromClaims, ok := token.Claims["xsrfToken"]
				if !ok {
					logger.Debug("The XSRF Token is missing from the JWT claims")
					http.Error(w, "Request unauthorized", http.StatusUnauthorized)
					return
				}

				xsrfTokenFromHeader := r.Header.Get("X-XSRF-TOKEN")

				if xsrfTokenFromHeader == "" || xsrfTokenFromClaims != xsrfTokenFromHeader {
					logger.Debug("The XSRF token does not match the XSRF claim")
					http.Error(w, "Request unauthorized", http.StatusUnauthorized)
					return
				}
			}
		} else {
			logger.Debugf("No AuthenticationToken cookie found: %s", err.Error())
		}

		// Verify the access token if no JWT was provided
		if err != nil {
			token, err = verifyAccessToken(r)
		}

		// If no JWT or access token found
		if err != nil {
			logger.Debug("No access token provided")
			http.Error(w, "Request unauthorized", http.StatusUnauthorized)
			return
		}

		setJWTInContext(r, token)
		next.ServeHTTP(w, r)
		context.Clear(r)
		return
	})
}

// Authenticate calls the proper handler based on whether authentication is enabled or not
func (c *Config) Authenticate(next http.Handler) http.Handler {
	if c.DriverName == "none" {
		return publicHandler(next)
	}
	return restrictedHandler(next)
}

// Login authenticates a user against the authentication driver
func (c *Config) Login() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
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

			user, err := c.login(u, p)
			if err != nil {
				logger.Info(err)

				// Output to audit log
				log := structs.AuditLog{Action: "loginfailure", Level: "default", Output: err.Error()}
				log.RemoteAddr = helpers.GetIP(r)
				audit.Log(log)

				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			xsrfToken := helpers.RandomString(32)
			authenticationToken, err := GetToken(user, xsrfToken)
			if err != nil {
				logger.Infof("Authentication failed, could not create the token: %s", err)
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			// Set the required cookies
			SetCookies(w, authenticationToken, xsrfToken)
			return
		}

		http.Redirect(w, r, "/#/login", http.StatusFound)
		return
	})
}
