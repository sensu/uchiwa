package config

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/palourde/mergo"
	"github.com/sensu/uchiwa/uchiwa/auth"
	"github.com/sensu/uchiwa/uchiwa/logger"
)

var (
	defaultGlobalConfig = GlobalConfig{
		Host:    "0.0.0.0",
		Port:    3000,
		Refresh: 10,
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
	}
	defaultSensuConfig = SensuConfig{
		Port:    4567,
		Timeout: 10,
	}
	defaultConfig = Config{
		Uchiwa: defaultGlobalConfig,
	}
)

// Load retrieves a specified configuration file and return a Config struct
func Load(file, directories string) *Config {
	// Load the configuration file
	conf, err := loadFile(file)
	if err != nil {
		logger.Fatal(err)
	}

	// Apply default configs to the configuration file
	if err := mergo.Merge(conf, defaultConfig); err != nil {
		logger.Fatal(err)
	}
	for i := range conf.Sensu {
		if err := mergo.Merge(&conf.Sensu[i], defaultSensuConfig); err != nil {
			logger.Fatal(err)
		}
	}

	if directories != "" {
		configDir := loadDirectories(directories)
		// Overwrite the file config with the configs from the directories
		if err := mergo.MergeWithOverwrite(conf, configDir); err != nil {
			logger.Fatal(err)
		}
	}

	conf.Sensu = initSensu(conf.Sensu)

	// Support the dashboard attribute
	if conf.Dashboard != nil {
		conf.Uchiwa = *conf.Dashboard
		// Apply the default config to the dashboard attribute
		if err := mergo.Merge(conf, defaultConfig); err != nil {
			logger.Fatal(err)
		}
	}

	conf.Uchiwa = initUchiwa(conf.Uchiwa)
	return conf
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
	logger.Infof("Loading the configuration file %s", path)

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
	} else if global.Gitlab.Server != "" {
		global.Auth.Driver = "gitlab"
	} else if global.Ldap.Server != "" {
		global.Auth.Driver = "ldap"
		if global.Ldap.GroupBaseDN == "" {
			global.Ldap.GroupBaseDN = global.Ldap.BaseDN
		}
		if global.Ldap.UserBaseDN == "" {
			global.Ldap.UserBaseDN = global.Ldap.BaseDN
		}
	} else if global.Db.Driver != "" && global.Db.Scheme != "" {
		global.Auth.Driver = "sql"
	} else if len(global.Users) != 0 {
		logger.Debug("Loading multiple users from the config")
		global.Auth.Driver = "simple"
	} else if global.User != "" && global.Pass != "" {
		logger.Debug("Loading single user from the config")
		global.Auth.Driver = "simple"

		// Support multiple users
		global.Users = append(global.Users, auth.User{Username: global.User, Password: global.Pass, FullName: global.User})
	}

	return global
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
	p.Uchiwa.Ldap.BindPass = "*****"
	p.Sensu = make([]SensuConfig, len(c.Sensu))
	for i := range c.Sensu {
		p.Sensu[i] = c.Sensu[i]
		p.Sensu[i].User = "*****"
		p.Sensu[i].Pass = "*****"
	}
	return p
}
