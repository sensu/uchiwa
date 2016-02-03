package config

import (
	"github.com/sensu/uchiwa/uchiwa/auth"
	"github.com/sensu/uchiwa/uchiwa/structs"
)

// Config struct contains []SensuConfig and UchiwaConfig structs
type Config struct {
	Dashboard *GlobalConfig `json:",omitempty"`
	Sensu     []SensuConfig
	Uchiwa    GlobalConfig
}

// SensuConfig struct contains conf about a Sensu API
type SensuConfig struct {
	Name     string
	Host     string
	Port     int
	Ssl      bool
	Insecure bool
	URL      string
	User     string
	Path     string
	Pass     string
	Timeout  int
}

// GlobalConfig struct contains conf about Uchiwa
type GlobalConfig struct {
	Host       string
	Port       int
	LogLevel   string
	Refresh    int
	Pass       string
	User       string
	Users      []auth.User
	Audit      Audit
	Auth       structs.Auth
	Db         Db
	Enterprise bool
	Github     Github
	Gitlab     Gitlab
	Ldap       Ldap
	SSL        SSL
}

// Audit struct contains the config of the Audit logger
type Audit struct {
	Level   string
	Logfile string
}

// Db struct contains the SQL driver configuration
type Db struct {
	Driver string
	Scheme string
}

// Github struct contains the GitHub driver configuration
type Github struct {
	ClientID     string
	ClientSecret string
	Roles        []auth.Role
	Server       string
}

// Gitlab struct contains the Gitlab driver configuration
type Gitlab struct {
	ApplicationID string
	Secret        string
	RedirectURL   string
	Roles         []auth.Role
	Server        string
}

// Ldap struct contains the LDAP driver configuration
type Ldap struct {
	Server               string
	Port                 int
	BaseDN               string
	BindUser             string
	BindPass             string
	GroupBaseDN          string
	GroupObjectClass     string
	GroupMemberAttribute string
	Insecure             bool
	Roles                []auth.Role
	Security             string
	UserAttribute        string
	UserBaseDN           string
	UserObjectClass      string
}

// SSL struct contains the path the SSL certificate and key
type SSL struct {
	CertFile string
	KeyFile  string
}
