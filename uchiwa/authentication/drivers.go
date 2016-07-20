package authentication

import "fmt"

// Advanced function allows a third party Identification driver
func (a *Config) Advanced(driver loginFn, driverName string) {
	a.DriverFn = driver
	a.DriverName = driverName

	initToken(a.Auth)
}

// None function sets the Config struct in order to disable authentication
func (a *Config) None() {
	a.DriverFn = none
	a.DriverName = "none"
}

// Simple function sets the Config struct in order to enable simple authentication based on provided user and pass
func (a *Config) Simple(u []User) {
	a.DriverFn = simple
	a.DriverName = "simple"

	users = u

	initToken(a.Auth)
}

// none represents the authentication driver when auth is disabled
func none(u, p string) (*User, error) {
	return &User{}, nil
}

// simple represents the simple authentication driver
func simple(u, p string) (*User, error) {
	for _, user := range users {
		if u == user.Username && p == user.Password {
			return &user, nil
		}
	}
	return &User{}, fmt.Errorf("invalid user '%s' or invalid password", u)
}
