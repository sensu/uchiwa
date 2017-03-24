package authentication

import (
	"fmt"
	"net/http"
	"time"
)

const (
	authenticationCookieName = "AuthenticationToken"
	xsrfCookieName           = "XSRF-TOKEN"
)

func (a *Config) login(user, pass string) (*User, error) {
	// Authenticate the user with the authentication driver
	u, err := a.DriverFn(user, pass)
	if err != nil {
		return nil, fmt.Errorf("Authentication failed: %s", err)
	}

	// Obfuscate the user's salt & hash
	u.PasswordHash = ""
	u.PasswordSalt = ""

	return u, nil
}

// DeleteCookies invalidate the JWT and XSTF cookies
func DeleteCookies(w http.ResponseWriter) {
	authenticationCookie := http.Cookie{
		Name:     authenticationCookieName,
		Value:    "",
		HttpOnly: true,
		Expires:  time.Now().Add(-100 * time.Hour),
		MaxAge:   -1,
	}
	http.SetCookie(w, &authenticationCookie)

	xsrfCookie := http.Cookie{
		Name:    xsrfCookieName,
		Value:   "",
		Expires: time.Now().Add(-100 * time.Hour),
		MaxAge:  -1,
	}
	http.SetCookie(w, &xsrfCookie)
}

// SetCookies set the proper cookies for the JWT and XSFR tokens
func SetCookies(w http.ResponseWriter, r *http.Request, authenticationToken, xsrfToken string) {
	var isSecure bool
	if r.TLS != nil {
		isSecure = true
	}

	authenticationCookie := http.Cookie{
		Name:     authenticationCookieName,
		Value:    authenticationToken,
		HttpOnly: true,
		Path:     "/",
		Secure:   isSecure,
	}
	http.SetCookie(w, &authenticationCookie)

	xsrfCookie := http.Cookie{
		Name:   xsrfCookieName,
		Value:  xsrfToken,
		Path:   "/",
		Secure: isSecure,
	}
	http.SetCookie(w, &xsrfCookie)
}
