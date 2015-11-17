package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	_, err := Load("../foo.bar")
	assert.NotNil(t, err, "should return an error when file does not exist")

	_, err = Load("../uchiwa.go")
	assert.NotNil(t, err, "should return an error when it cannot parse a file")

	conf, err := Load("../../fixtures/config_test.json")
	assert.Nil(t, err, "got unexpected error: %s", err)

	// private config
	assert.NotEqual(t, "*****", conf.Uchiwa.User, "Uchiwa user in private config shouldn't be masked")
	assert.NotEqual(t, "*****", conf.Uchiwa.Pass, "Uchiwa pass in private config shouldn't be masked")
	for i := range conf.Sensu {
		assert.NotEqual(t, "*****", conf.Sensu[i].User, "Sensu APIs user in private config shouldn't be masked")
		assert.NotEqual(t, "*****", conf.Sensu[i].Pass, "Sensu APIs pass in private config shouldn't be masked")
	}
	assert.Equal(t, 1, len(conf.Uchiwa.Users))

	// public config
	public := conf.GetPublic()
	assert.Equal(t, "*****", public.Uchiwa.User, "Uchiwa user in public config should be masked")
	assert.Equal(t, "*****", public.Uchiwa.Pass, "Uchiwa pass in public config should be masked")
	for i := range public.Sensu {
		assert.Equal(t, "*****", public.Sensu[i].User, "Sensu APIs user in public config should be masked")
		assert.Equal(t, "*****", public.Sensu[i].Pass, "Sensu APIs pass in public config should be masked")
	}
	assert.Equal(t, 0, len(public.Uchiwa.Users))

}

func TestLoadArrayOfUsers(t *testing.T) {
	conf, err := Load("../../fixtures/config_test_multiple.json")
	assert.Nil(t, err, "got unexpected error: %s", err)
	assert.NotNil(t, conf, "conf should not be nil")

	assert.Equal(t, "simple", conf.Uchiwa.Auth.Driver, "Uchiwa authentication driver should be 'simple'")
	assert.Equal(t, 2, len(conf.Uchiwa.Users))
}

func TestLoadArrayOfUsersOnPublicGet(t *testing.T) {
	conf, err := Load("../../fixtures/config_test_multiple.json")
	assert.Nil(t, err, "got unexpected error: %s", err)
	assert.NotNil(t, conf, "conf should not be nil")

	assert.Equal(t, "simple", conf.Uchiwa.Auth.Driver, "Uchiwa authentication driver should be 'simple'")
	public := conf.GetPublic()
	assert.Equal(t, 0, len(public.Uchiwa.Users))
}

func TestInitSensu(t *testing.T) {
	c := Config{
		Sensu: []SensuConfig{
			{Name: "foo ? bar", Host: "127.0.0.1"},
			{Name: "bar / foo", Host: "127.0.0.1"},
		},
	}

	c.initSensu()

	expectedConfig := []SensuConfig{
		{Name: "foo  bar", Host: "127.0.0.1", Port: 4567, Ssl: false, Insecure: false, URL: "http://127.0.0.1:4567", User: "", Path: "", Pass: "", Timeout: 10},
		{Name: "bar  foo", Host: "127.0.0.1", Port: 4567, Ssl: false, Insecure: false, URL: "http://127.0.0.1:4567", User: "", Path: "", Pass: "", Timeout: 10},
	}
	assert.Equal(t, expectedConfig, c.Sensu)
}
