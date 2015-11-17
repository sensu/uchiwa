package auth

import "github.com/sensu/uchiwa/uchiwa/structs"

// Config struct contains the authentication configuration
type Config struct {
	Auth       structs.Auth
	DriverFn   loginFn
	DriverName string
}

// User structure
type User struct {
	ID           int64
	Username     string
	FullName     string
	Email        string
	Password     string
	PasswordHash string
	PasswordSalt string
	Role         Role
	Token        string
}

type loginFn func(string, string) (*User, error)

var (
	users []User
)

// New function initalizes and returns a Config struct
func New(auth structs.Auth) Config {
	a := Config{
		Auth: auth,
	}
	return a
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

// Advanced function allows a third party Identification driver
func (a *Config) Advanced(driver loginFn, driverName string) {
	a.DriverFn = driver
	a.DriverName = driverName

	initToken(a.Auth)
}
