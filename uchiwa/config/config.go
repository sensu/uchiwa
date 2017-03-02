package config

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/palourde/mergo"
	"github.com/sensu/uchiwa/uchiwa/authentication"
	"github.com/sensu/uchiwa/uchiwa/logger"
)

const obfuscatedValue = "*****"

var (
	defaultGlobalConfig = GlobalConfig{
		Host:     "0.0.0.0",
		Port:     3000,
		LogLevel: "info",
		Refresh:  10,
		Ldap: Ldap{
			Port:                 389,
			Security:             "none",
			UserAttribute:        "sAMAccountName",
			UserObjectClass:      "person",
			GroupMemberAttribute: "member",
			GroupObjectClass:     "groupOfNames",
		},
		Audit: Audit{
			Level:   "default",
			Logfile: "/var/log/sensu/sensu-enterprise-dashboard-audit.log",
		},
		UsersOptions: UsersOptions{
			DateFormat:             "YYYY-MM-DD HH:mm:ss",
			DefaultExpireOnResolve: false,
			DefaultTheme:           "uchiwa-default",
			DisableNoExpiration:    false,
			RequireSilencingReason: false,
		},
	}
	defaultSensuConfig = SensuConfig{
		Port:    4567,
		Timeout: 10,
	}
	defaultConfig = Config{
		Uchiwa: defaultGlobalConfig,
	}
	// Private contains the private configuration
	Private *Config
)

// Load retrieves the Uchiwa configuration from files and directories
// and returns the private configuration as a Config struct pointer
func Load(file, directories string) *Config {
	// Load the configuration file
	var err error
	Private, err = loadFile(file)
	if err != nil {
		logger.Fatal(err)
	}

	// Apply default configs to the configuration file
	if err := mergo.Merge(Private, defaultConfig); err != nil {
		logger.Fatal(err)
	}
	for i := range Private.Sensu {
		if err := mergo.Merge(&Private.Sensu[i], defaultSensuConfig); err != nil {
			logger.Fatal(err)
		}
	}

	if directories != "" {
		configDir := loadDirectories(directories)
		// Overwrite the file config with the configs from the directories
		if err := mergo.MergeWithOverwrite(Private, configDir); err != nil {
			logger.Fatal(err)
		}
	}

	Private.Sensu = initSensu(Private.Sensu)

	// Support the dashboard attribute
	if Private.Dashboard != nil {
		Private.Uchiwa = *Private.Dashboard
		// Apply the default config to the dashboard attribute
		if err := mergo.Merge(Private, defaultConfig); err != nil {
			logger.Fatal(err)
		}
	}

	Private.Uchiwa = initUchiwa(Private.Uchiwa)
	return Private
}

// loadDirectories loads a Config struct from one or multiple directories of configuration
func loadDirectories(path string) *Config {
	conf := new(Config)
	var configFiles []string
	directories := strings.Split(strings.ToLower(path), ",")

	for _, directory := range directories {
		// Find all JSON files in the specified directories
		files, err := filepath.Glob(filepath.Join(directory, "*.json"))
		if err != nil {
			logger.Warning(err)
			continue
		}

		// Add the files found to a slice of configuration files to open
		for _, file := range files {
			configFiles = append(configFiles, file)
		}
	}

	// Load every configuration files and merge them together bit by bit
	for _, file := range configFiles {
		// Load the config from the file
		c, err := loadFile(file)
		if err != nil {
			logger.Warning(err)
			continue
		}

		// Apply this configuration to the existing one
		if err := mergo.MergeWithOverwrite(conf, c); err != nil {
			logger.Warning(err)
			continue
		}
	}

	// Apply the default config to the Sensu APIs
	for i := range conf.Sensu {
		if err := mergo.Merge(&conf.Sensu[i], defaultSensuConfig); err != nil {
			logger.Fatal(err)
		}
	}

	return conf
}

// loadFile loads a Config struct from a configuration file
func loadFile(path string) (*Config, error) {
	logger.Warningf("Loading the configuration file %s", path)

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

	return c, nil
}

func initSensu(apis []SensuConfig) []SensuConfig {
	for i, api := range apis {
		// Set a datacenter name if missing
		if api.Name == "" {
			logger.Warningf("Sensu API %s has no name property, make sure to set it in your configuration. Generating a temporary one...", api.URL)
			apis[i].Name = fmt.Sprintf("sensu-%v", rand.Intn(100))
		}

		// Escape special characters in DC name
		r := strings.NewReplacer(":", "", "/", "", ";", "", "?", "")
		apis[i].Name = r.Replace(apis[i].Name)

		// Make sure the host is not empty
		if api.Host == "" {
			logger.Fatalf("Sensu API %q Host is missing", api.Name)
		}

		// Determine the protocol to use
		prot := "http"
		if api.Ssl {
			prot += "s"
		}

		// Set the API URL
		apis[i].URL = fmt.Sprintf("%s://%s:%d%s", prot, api.Host, api.Port, api.Path)
	}
	return apis
}

func initUchiwa(global GlobalConfig) GlobalConfig {

	// Set the proper authentication driver
	if global.Github.Server != "" {
		global.Auth.Driver = "github"

		for i := range global.Github.Roles {
			authentication.Roles = append(authentication.Roles, global.Github.Roles[i])
		}
	} else if global.Gitlab.Server != "" {
		global.Auth.Driver = "gitlab"

		for i := range global.Gitlab.Roles {
			authentication.Roles = append(authentication.Roles, global.Gitlab.Roles[i])
		}
	} else if global.Ldap.Server != "" {
		global.Auth.Driver = "ldap"
		if global.Ldap.GroupBaseDN == "" {
			global.Ldap.GroupBaseDN = global.Ldap.BaseDN
		}
		if global.Ldap.UserBaseDN == "" {
			global.Ldap.UserBaseDN = global.Ldap.BaseDN
		}

		for i := range global.Ldap.Roles {
			authentication.Roles = append(authentication.Roles, global.Ldap.Roles[i])
		}
	} else if global.OIDC.Server != "" {
		global.Auth.Driver = "oidc"

		for i := range global.OIDC.Roles {
			authentication.Roles = append(authentication.Roles, global.OIDC.Roles[i])
		}
	} else if global.Db.Driver != "" && global.Db.Scheme != "" {
		global.Auth.Driver = "sql"
	} else if len(global.Users) != 0 {
		logger.Debug("Loading multiple users from the config")
		global.Auth.Driver = "simple"

		for i := range global.Users {
			if global.Users[i].AccessToken != "" {
				global.Users[i].Role.AccessToken = global.Users[i].AccessToken
			}
			if global.Users[i].Readonly != false {
				global.Users[i].Role.Readonly = global.Users[i].Readonly
			}
			authentication.Roles = append(authentication.Roles, global.Users[i].Role)
		}
	} else if global.User != "" && global.Pass != "" {
		logger.Debug("Loading single user from the config")
		global.Auth.Driver = "simple"

		// Support multiple users
		global.Users = append(global.Users, authentication.User{Username: global.User, Password: global.Pass, FullName: global.User})
	}

	// Set the logger level
	logger.SetLogLevel(global.LogLevel)

	// Set the refresh rate for frontend
	global.UsersOptions.Refresh = global.Refresh * 1000

	return global
}

// GetPublic generates the public configuration
func (c *Config) GetPublic() *Config {
	p := new(Config)
	p.Uchiwa = c.Uchiwa
	p.Uchiwa.User = obfuscatedValue
	p.Uchiwa.Pass = obfuscatedValue
	p.Uchiwa.Users = []authentication.User{}
	p.Uchiwa.Db.Scheme = obfuscatedValue
	p.Uchiwa.Github.ClientID = obfuscatedValue
	p.Uchiwa.Github.ClientSecret = obfuscatedValue
	p.Uchiwa.Gitlab.ClientID = obfuscatedValue
	p.Uchiwa.Gitlab.ClientSecret = obfuscatedValue
	p.Uchiwa.Ldap.BindPass = obfuscatedValue
	p.Uchiwa.OIDC.ClientID = obfuscatedValue
	p.Uchiwa.OIDC.ClientSecret = obfuscatedValue

	for i := range p.Uchiwa.Github.Roles {
		p.Uchiwa.Github.Roles[i].AccessToken = obfuscatedValue
	}

	for i := range p.Uchiwa.Gitlab.Roles {
		p.Uchiwa.Gitlab.Roles[i].AccessToken = obfuscatedValue
	}

	for i := range p.Uchiwa.Ldap.Roles {
		p.Uchiwa.Ldap.Roles[i].AccessToken = obfuscatedValue
	}

	for i := range p.Uchiwa.OIDC.Roles {
		p.Uchiwa.OIDC.Roles[i].AccessToken = obfuscatedValue
	}

	p.Sensu = make([]SensuConfig, len(c.Sensu))
	for i := range c.Sensu {
		p.Sensu[i] = c.Sensu[i]
		p.Sensu[i].User = obfuscatedValue
		p.Sensu[i].Pass = obfuscatedValue
	}

	return p
}
