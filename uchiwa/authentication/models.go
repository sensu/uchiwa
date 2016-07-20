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
	Name          string
	Readonly      bool
	Scope         Scope
	Subscriptions []string
}

// Scope contains the type of access of a role
type Scope struct {
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
