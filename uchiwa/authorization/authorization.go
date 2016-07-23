package authorization

import (
	"encoding/json"
	"net/http"

	"github.com/sensu/uchiwa/uchiwa/authentication"
	"github.com/sensu/uchiwa/uchiwa/logger"
)

// Authorization contains the different methods used for authorizing
// requests made to the API
type Authorization interface {
	Handler(http.Handler) http.Handler
}

// Uchiwa represents an instance of the Authorization interface for the community
type Uchiwa struct{}

// Handler verifies if the user has access to the requested resource
func (u *Uchiwa) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		readonly := isReadOnly(r)
		authorized := isAuthorized(readonly, r.Method)
		if !authorized {
			http.Error(w, "Request forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAuthorized(isReadOnly bool, method string) bool {
	if (method != "HEAD" && method != "GET") && isReadOnly {
		return false
	}
	return true
}

// hasReadOnly verifies if the user only has read-only access.
// Returns true if the user only have read-only access
func isReadOnly(r *http.Request) bool {
	var role authentication.Role

	token := authentication.GetJWTFromContext(r)
	if token == nil { // authentication is not enabled
		logger.Debug("No JWT found in context")
		return false
	}

	m := token.Claims["Role"]

	// use JSON representation of the interface to assert it into the uchiwa.Role struct
	j, err := json.Marshal(&m)
	if err != nil {
		// Unexpected error, fallback to readonly user
		logger.Warning("Could not marshal the token claims")
		return true
	}

	err = json.Unmarshal(j, &role)
	if err != nil {
		// Unexpected error, fallback to readonly user
		logger.Warning("Could not unmarshal the token claims")
		return true
	}

	return role.Readonly
}
