package authentication

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sensu/uchiwa/uchiwa/logger"
)

// TokenLocation represents a function that accepts a request as input and returns
// either a token or an error.
type TokenLocation func(r *http.Request) (string, error)

// accessTokenFromAuthHeader retrieves an access token from the
// authentication header, e.g. token {TOKEN}
func accessTokenFromAuthHeader(r *http.Request) (string, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return "", nil // No token found, do not fail yet
	}

	authorizationComponents := strings.Split(authorization, " ")
	if len(authorizationComponents) != 2 || strings.ToLower(authorizationComponents[0]) != "token" {
		logger.Debug("Invalid authorization header. The format must be: token {token}")
		return "", errors.New("")
	}

	return authorizationComponents[1], nil
}

// accessTokenFromParameter retrieves an access token from a
// URL parameter, e.g. ?token={TOKEN}
func accessTokenFromParameter(r *http.Request) (string, error) {
	return r.URL.Query().Get("token"), nil
}

// findAccessToken finds the first corresponding token in the provided functions
func findAccessToken(locations ...TokenLocation) TokenLocation {
	return func(r *http.Request) (string, error) {
		for _, location := range locations {
			token, err := location(r)
			if err != nil {
				return "", err
			}
			if token != "" {
				return token, nil
			}
		}
		return "", errors.New("")
	}
}

// findRoleFromAccessToken finds within the Role slice a role with
// the corresponding token
func findRoleFromAccessToken(token string) (*Role, error) {
	for _, role := range Roles {
		if role.AccessToken == token {
			return &role, nil
		}
	}
	return nil, fmt.Errorf("No role found with the access token %s", token)
}

// verifyAccessToken extracts the access token and verifies it, then
// returns a JWT with the associated role
func verifyAccessToken(r *http.Request) (*jwt.Token, error) {
	locations := findAccessToken(accessTokenFromAuthHeader, accessTokenFromParameter)
	accessToken, err := locations(r)
	if err != nil {
		return nil, errors.New("")
	}

	role, err := findRoleFromAccessToken(accessToken)
	if err != nil {
		return nil, errors.New("")
	}

	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["role"] = role
	token.Claims["username"] = role.Name

	return token, nil
}
