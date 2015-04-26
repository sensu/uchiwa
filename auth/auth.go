package auth

// Config struct contains the authentication configuration
type Config struct {
	Driver     loginFn
	DriverName string
}

// User structure
type User struct {
	ID           int64  `db:"id"`
	Username     string `db:"username"`
	FullName     string `db:"fullName"`
	Email        string `db:"email"`
	Password     string `db:"-"`
	PasswordHash string `db:"passwordHash"`
	PasswordSalt string `db:"passwordSalt"`
	Role         string `db:"role"`
	Token        string `db:"-"`
}

type loginFn func(string, string) (*User, error)

var (
	user, pass string
)

// New function initalizes and returns a Config struct
func New() Config {
	a := Config{}
	return a
}

// None function sets the Config struct in order to disable authentication
func (a *Config) None() {
	a.Driver = none
	a.DriverName = "none"
}

// Simple function sets the Config struct in order to enable simple authentication based on provided user and pass
func (a *Config) Simple(u, p string) {
	a.Driver = simple
	a.DriverName = "simple"

	user = u
	pass = p

	initToken()
}

// Advanced function allows a third party Identification driver
func (a *Config) Advanced(driver loginFn, driverName string) {
	a.Driver = driver
	a.DriverName = driverName

	initToken()
}
