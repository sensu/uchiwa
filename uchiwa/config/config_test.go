package config

import (
	"testing"

	"github.com/sensu/uchiwa/uchiwa/authentication"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Only default & config file
	conf := Load("../../fixtures/config_test.json", "")
	assert.Equal(t, 4567, conf.Sensu[0].Port)
	assert.Equal(t, 10, conf.Sensu[0].Timeout)
	assert.Equal(t, 10, conf.Uchiwa.Refresh)
	assert.Equal(t, false, conf.Uchiwa.DisableNoExpiration)
	assert.Equal(t, false, conf.Uchiwa.ExpireOnResolveDefault)
	assert.Equal(t, 389, conf.Uchiwa.Ldap.Port)
	assert.Equal(t, "person", conf.Uchiwa.Ldap.UserObjectClass)
	assert.Equal(t, "default", conf.Uchiwa.Audit.Level)

	conf = Load("../../fixtures/config_test.json", "../../fixtures/conf.d")
	assert.Equal(t, 5, len(conf.Sensu))
	assert.Equal(t, 4567, conf.Sensu[0].Port)
	assert.Equal(t, 10, conf.Sensu[0].Timeout)
	assert.Equal(t, 4569, conf.Sensu[3].Port)
	assert.Equal(t, 5, conf.Sensu[3].Timeout)
	assert.Equal(t, 4567, conf.Sensu[4].Port)
	assert.Equal(t, 10, conf.Sensu[4].Timeout)

	assert.Equal(t, "192.168.0.1", conf.Uchiwa.Host)
	assert.Equal(t, 8000, conf.Uchiwa.Port)
	assert.Equal(t, 2, len(conf.Uchiwa.Users))

	// Test the removal of the Role object in configuration files
	assert.Equal(t, false, conf.Uchiwa.Users[0].Role.Readonly)
	assert.Equal(t, "foobar", conf.Uchiwa.Users[0].Role.AccessToken)
	assert.Equal(t, true, conf.Uchiwa.Users[1].Role.Readonly)

	// We should also support the dashboard attribute instead of uchiwa
	conf = Load("../../fixtures/config_dashboard.json", "")
	assert.Equal(t, "127.0.0.1", conf.Uchiwa.Host)
	assert.Equal(t, 8080, conf.Uchiwa.Port)
	assert.Equal(t, 1, len(conf.Uchiwa.Users))
	assert.Equal(t, 389, conf.Uchiwa.Ldap.Port)
	assert.Equal(t, "person", conf.Uchiwa.Ldap.UserObjectClass)
}

func TestLoadDirectories(t *testing.T) {
	conf := loadDirectories("foobar,../../fixtures/conf.d")

	assert.Equal(t, 3, len(conf.Sensu))
	assert.Equal(t, "us-east-2", conf.Sensu[0].Name)
	assert.Equal(t, "us-west-2", conf.Sensu[1].Name)
	assert.Equal(t, "us-east-3", conf.Sensu[2].Name)

	assert.Equal(t, 2, len(conf.Uchiwa.Users))
	assert.Equal(t, "admin", conf.Uchiwa.Users[0].Username)
	assert.Equal(t, "readonly", conf.Uchiwa.Users[1].Username)
	assert.Equal(t, "192.168.0.1", conf.Uchiwa.Host)
	assert.Equal(t, 8000, conf.Uchiwa.Port)
}

func TestLoadFile(t *testing.T) {
	_, err := loadFile("foo.bar")
	assert.NotNil(t, err, "foo.bar does not exist")

	_, err = loadFile("config.go")
	assert.NotNil(t, err, "config.go is not a JSON file")

	conf, err := loadFile("../../fixtures/config_test.json")
	assert.Nil(t, err, "got unexpected error: %s", err)

	// Sensu APIs
	assert.Equal(t, 2, len(conf.Sensu))
	assert.Equal(t, "us-east-1", conf.Sensu[0].Name)
	assert.Equal(t, "us-west-1", conf.Sensu[1].Name)
	assert.Equal(t, 4570, conf.Sensu[1].Port)
	assert.Equal(t, 5, conf.Sensu[1].Timeout)

	// Uchiwa
	assert.Equal(t, "0.0.0.0", conf.Uchiwa.Host)
	assert.Equal(t, 8080, conf.Uchiwa.Port)
	assert.Equal(t, "foo", conf.Uchiwa.User)
	assert.Equal(t, "bar", conf.Uchiwa.Pass)
}

func TestInitSensu(t *testing.T) {
	apis := []SensuConfig{
		SensuConfig{Host: "10.0.0.1", Port: 4567},
		SensuConfig{Name: "test/1", Host: "10.0.10.1", Port: 4567, Ssl: true},
	}

	sensu := initSensu(apis)
	assert.NotEqual(t, "", sensu[0].Name)
	assert.Equal(t, "http://10.0.0.1:4567", sensu[0].URL)
	assert.Equal(t, "test1", sensu[1].Name)
	assert.Equal(t, "https://10.0.10.1:4567", sensu[1].URL)
}

func TestInitUchiwa(t *testing.T) {
	conf := GlobalConfig{Github: Github{Server: "127.0.0.1"}}
	uchiwa := initUchiwa(conf)
	assert.Equal(t, "github", uchiwa.Auth.Driver)

	conf = GlobalConfig{Gitlab: Gitlab{Server: "127.0.0.1"}}
	uchiwa = initUchiwa(conf)
	assert.Equal(t, "gitlab", uchiwa.Auth.Driver)

	conf = GlobalConfig{Ldap: Ldap{BaseDN: "cn=foo", Server: "127.0.0.1"}}
	uchiwa = initUchiwa(conf)
	assert.Equal(t, "ldap", uchiwa.Auth.Driver)
	assert.Equal(t, Ldap{BaseDN: "cn=foo", GroupBaseDN: "cn=foo", UserBaseDN: "cn=foo", Server: "127.0.0.1"}, uchiwa.Ldap)

	conf = GlobalConfig{Db: Db{Driver: "mysql", Scheme: "foo"}}
	uchiwa = initUchiwa(conf)
	assert.Equal(t, "sql", uchiwa.Auth.Driver)

	conf = GlobalConfig{Users: []authentication.User{authentication.User{ID: 1}}}
	uchiwa = initUchiwa(conf)
	assert.Equal(t, "simple", uchiwa.Auth.Driver)

	conf = GlobalConfig{User: "foo", Pass: "secret"}
	uchiwa = initUchiwa(conf)
	assert.Equal(t, "simple", uchiwa.Auth.Driver)
	assert.Equal(t, []authentication.User{authentication.User{ID: 0, FullName: "foo", Password: "secret", Username: "foo"}}, uchiwa.Users)
}

func TestGetPublic(t *testing.T) {
	conf := Config{
		Sensu: []SensuConfig{
			SensuConfig{
				User: "foo",
				Pass: "secret",
			},
		},
		Uchiwa: GlobalConfig{
			User:   "foo",
			Pass:   "secret",
			Users:  []authentication.User{authentication.User{ID: 1}},
			Db:     Db{Scheme: "foo"},
			Github: Github{ClientID: "foo", ClientSecret: "secret"},
			Ldap:   Ldap{BindPass: "secret"},
		},
	}

	pubConf := conf.GetPublic()

	assert.NotEqual(t, conf, pubConf)

	assert.Equal(t, "foo", conf.Sensu[0].User)
	assert.Equal(t, "secret", conf.Sensu[0].Pass)
	assert.Equal(t, "foo", conf.Uchiwa.User)
	assert.Equal(t, "secret", conf.Uchiwa.Pass)

	assert.Equal(t, "*****", pubConf.Sensu[0].User)
	assert.Equal(t, "*****", pubConf.Sensu[0].Pass)
	assert.Equal(t, "*****", pubConf.Uchiwa.User)
	assert.Equal(t, "*****", pubConf.Uchiwa.Pass)
	assert.Equal(t, []authentication.User{}, pubConf.Uchiwa.Users)
	assert.Equal(t, "*****", pubConf.Uchiwa.Db.Scheme)
	assert.Equal(t, "*****", pubConf.Uchiwa.Github.ClientID)
	assert.Equal(t, "*****", pubConf.Uchiwa.Github.ClientSecret)
	assert.Equal(t, "*****", pubConf.Uchiwa.Ldap.BindPass)
}
