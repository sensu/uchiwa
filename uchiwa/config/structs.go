package config

import (
	"crypto/tls"

	"github.com/sensu/uchiwa/uchiwa/authentication"
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
	Host         string
	Port         int
	LogLevel     string
	Refresh      int
	Pass         string
	User         string
	Users        []authentication.User
	Audit        Audit
	Auth         structs.Auth
	Db           Db
	Enterprise   bool
	Github       Github
	Gitlab       Gitlab
	Ldap         Ldap
	OIDC         OIDC
	SSL          SSL
	UsersOptions UsersOptions
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
	Roles        []authentication.Role
	Server       string
}

// Gitlab struct contains the Gitlab driver configuration
type Gitlab struct {
	ClientID     string `json:"applicationid"`
	ClientSecret string `json:"secret"`
	RedirectURL  string
	Roles        []authentication.Role
	Server       string
}

// Ldap struct contains the LDAP driver configuration
type Ldap struct {
	LdapServer
	Debug   bool
	Roles   []authentication.Role
	Servers []LdapServer
}

type LdapServer struct {
	Server               string
	Port                 int
	BaseDN               string
	BindUser             string
	BindPass             string
	Dialect              string
	DisableNestedGroups  bool
	GroupBaseDN          string
	GroupObjectClass     string
	GroupMemberAttribute string
	Insecure             bool
	Security             string
	TLSConfig            *tls.Config
	UserAttribute        string
	UserBaseDN           string
	UserObjectClass      string
}

// OIDC struct contains the OIDC driver configuration
type OIDC struct {
	ClientID     string
	ClientSecret string
	Insecure     bool
	RedirectURL  string
	Roles        []authentication.Role
	Server       string
}

// SSL struct contains the path the SSL certificate and key
type SSL struct {
	CertFile string
	KeyFile  string
}

// UsersOptions struct contains various config tweaks
type UsersOptions struct {
	DateFormat             string
	DefaultTheme           string
	DisableNoExpiration    bool
	Favicon                string
	LogoURL                string
	Refresh                int
	RequireSilencingReason bool
	SilenceDurations       []float32
}
