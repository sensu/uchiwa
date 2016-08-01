package authentication

import "github.com/sensu/uchiwa/uchiwa/structs"

var (
	// Roles contains the roles for the active auth driver
	Roles []Role
	users []User
)

type loginFn func(string, string) (*User, error)

// Config contains the authentication configuration
type Config struct {
	Auth       structs.Auth
	DriverFn   loginFn
	DriverName string
}

// Role contains the attributes of a role
type Role struct {
	AccessToken   string
	Datacenters   []string
	Fallback      bool
	Members       []string
	Methods       Methods
	Name          string
	Readonly      bool
	Subscriptions []string
}

// Methods contains the allowed endpoints for each HTTP method
type Methods struct {
	Delete []string
	Head   []string
	Get    []string
	Post   []string
}

// User contains the attributes of a user
type User struct {
	ID           int64
	AccessToken  string
	Email        string
	FullName     string
	Password     string
	PasswordHash string
	PasswordSalt string
	Readonly     bool
	Role         Role
	Token        string
	Username     string
}
