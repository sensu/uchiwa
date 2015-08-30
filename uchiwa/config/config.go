package config

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/palourde/logger"
	"github.com/sensu/uchiwa/uchiwa/auth"
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
	Host    string
	Port    int
	Refresh int
	Pass    string
	User    string
	Users   []auth.User
	Db      Db
	Github  Github
	Ldap    Ldap
	Audit   Audit
	Auth    string
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

// Ldap struct contains the LDAP driver configuration
type Ldap struct {
	Server        string
	Port          int
	BaseDN        string
	BindUser      string
	BindPass      string
	Insecure      bool
	Roles         []auth.Role
	Security      string
	UserAttribute string
}

// Load retrieves a specified configuration file and return a Config struct
func Load(path string) (*Config, error) {
	logger.Infof("Loading configuration file %s", path)
	c := new(Config)
	file, err := os.Open(path)
	if err != nil {
		if len(path) > 1 {
			return nil, fmt.Errorf("Error: could not read config file %s.", path)
		}
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		return nil, fmt.Errorf("Error decoding file %s: %s", path, err)
	}

	c.initUchiwa()
	c.initSensu()

	return c, nil
}

func (c *Config) initSensu() {
	for i, api := range c.Sensu {
		prot := "http"
		if api.Name == "" {
			logger.Warningf("Sensu API %s has no name property. Generating random one...", api.URL)
			c.Sensu[i].Name = fmt.Sprintf("sensu-%v", rand.Intn(100))
		}
		// escape special characters in DC name
		r := strings.NewReplacer(":", "", "/", "", ";", "", "?", "")
		c.Sensu[i].Name = r.Replace(api.Name)

		if api.Host == "" {
			logger.Fatalf("Sensu API %q Host is missing", api.Name)
		}
		if api.Timeout == 0 {
			c.Sensu[i].Timeout = 10
		} else if api.Timeout >= 1000 { // backward compatibility with < 0.3.0 version
			c.Sensu[i].Timeout = api.Timeout / 1000
		}
		if api.Port == 0 {
			c.Sensu[i].Port = 4567
		}
		if api.Ssl {
			prot += "s"
		}
		c.Sensu[i].URL = fmt.Sprintf("%s://%s:%d%s", prot, api.Host, c.Sensu[i].Port, api.Path)
	}
}

func (c *Config) initUchiwa() {
	if c.Dashboard != nil {
		c.Uchiwa = *c.Dashboard
	}
	if c.Uchiwa.Host == "" {
		c.Uchiwa.Host = "0.0.0.0"
	}
	if c.Uchiwa.Port == 0 {
		c.Uchiwa.Port = 3000
	}
	if c.Uchiwa.Refresh == 0 {
		c.Uchiwa.Refresh = 10
	} else if c.Uchiwa.Refresh >= 1000 { // backward compatibility with < 0.3.0 version
		c.Uchiwa.Refresh = c.Uchiwa.Refresh / 1000
	}

	// authentication
	if c.Uchiwa.Github.Server != "" {
		c.Uchiwa.Auth = "github"
	} else if c.Uchiwa.Ldap.Server != "" {
		c.Uchiwa.Auth = "ldap"
		if c.Uchiwa.Ldap.Port == 0 {
			c.Uchiwa.Ldap.Port = 389
		}
		if c.Uchiwa.Ldap.Security == "" {
			c.Uchiwa.Ldap.Security = "none"
		}
		if c.Uchiwa.Ldap.UserAttribute == "" {
			c.Uchiwa.Ldap.UserAttribute = "sAMAccountName"
		}
	} else if c.Uchiwa.Db.Driver != "" && c.Uchiwa.Db.Scheme != "" {
		c.Uchiwa.Auth = "sql"
	} else if len(c.Uchiwa.Users) != 0 {
		logger.Debug("Loading multiple users from the config")
		c.Uchiwa.Auth = "simple"
	} else if c.Uchiwa.User != "" && c.Uchiwa.Pass != "" {
		logger.Debug("Loading single user from the config")
		c.Uchiwa.Auth = "simple"
		c.Uchiwa.Users = append(c.Uchiwa.Users, auth.User{Username: c.Uchiwa.User, Password: c.Uchiwa.Pass, FullName: c.Uchiwa.User})
	}

	// audit
	if c.Uchiwa.Audit.Level != "verbose" && c.Uchiwa.Audit.Level != "disabled" {
		c.Uchiwa.Audit.Level = "default"
	}
	if c.Uchiwa.Audit.Logfile == "" {
		c.Uchiwa.Audit.Logfile = "/var/log/sensu/sensu-enterprise-dashboard-audit.log"
	}
}

// GetPublic generates the public configuration
func (c *Config) GetPublic() *Config {
	p := new(Config)
	p.Uchiwa = c.Uchiwa
	p.Uchiwa.User = "*****"
	p.Uchiwa.Pass = "*****"
	p.Uchiwa.Users = []auth.User{}
	p.Uchiwa.Db.Scheme = "*****"
	p.Uchiwa.Github.ClientID = "*****"
	p.Uchiwa.Github.ClientSecret = "*****"
	p.Sensu = make([]SensuConfig, len(c.Sensu))
	for i := range c.Sensu {
		p.Sensu[i] = c.Sensu[i]
		p.Sensu[i].User = "*****"
		p.Sensu[i].Pass = "*****"
	}
	return p
}
