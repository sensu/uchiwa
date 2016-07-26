package authorization

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/sensu/uchiwa/uchiwa/authentication"
	"github.com/stretchr/testify/assert"
)

type HandleTester func(method string, params url.Values) *httptest.ResponseRecorder

func generateToken(role authentication.Role) *jwt.Token {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["Role"] = role

	return token
}

func mockHandlerRequest(method string, readonly bool) *httptest.ResponseRecorder {
	role := authentication.Role{Readonly: readonly}
	token := generateToken(role)

	r, _ := http.NewRequest(method, "/", nil)
	w := httptest.NewRecorder()

	setJWTInContext(r, token)

	handler := u.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	handler.ServeHTTP(w, r)

	return w
}

func setJWTInContext(r *http.Request, token *jwt.Token) {
	context.Set(r, "jwtToken", token)
}

var u *Uchiwa

func TestHandler(t *testing.T) {
	// GET with read-write
	w := mockHandlerRequest("GET", false)
	assert.Equal(t, http.StatusOK, w.Code)

	// GET with read-only
	w = mockHandlerRequest("GET", true)
	assert.Equal(t, http.StatusOK, w.Code)

	// POST with read-write
	w = mockHandlerRequest("POST", false)
	assert.Equal(t, http.StatusOK, w.Code)

	// POST with read-write
	w = mockHandlerRequest("POST", true)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestIsAuthorized(t *testing.T) {
	authorized := isAuthorized(false, "GET")
	assert.True(t, authorized)

	authorized = isAuthorized(true, "GET")
	assert.True(t, authorized)

	authorized = isAuthorized(true, "HEAD")
	assert.True(t, authorized)

	authorized = isAuthorized(false, "POST")
	assert.True(t, authorized)

	authorized = isAuthorized(true, "POST")
	assert.False(t, authorized)
}

func TestIsReadOnly(t *testing.T) {
	// Not readonly
	role := authentication.Role{Readonly: false}
	token := generateToken(role)
	r, _ := http.NewRequest("GET", "/", nil)
	setJWTInContext(r, token)
	readonly := isReadOnly(r)
	assert.False(t, readonly)

	// Is readonly
	role = authentication.Role{Readonly: true}
	token = generateToken(role)
	r, _ = http.NewRequest("GET", "/", nil)
	setJWTInContext(r, token)
	readonly = isReadOnly(r)
	assert.True(t, readonly)

	// Has no token
	r, _ = http.NewRequest("GET", "/", nil)
	readonly = isReadOnly(r)
	assert.False(t, readonly)
}
