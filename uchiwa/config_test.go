package uchiwa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	_, err := LoadConfig("../foo.bar")
	assert.NotNil(t, err, "should return an error when file does not exist")

	_, err = LoadConfig("../uchiwa.go")
	assert.NotNil(t, err, "should return an error when it cannot parse a file")

	conf, err := LoadConfig("../test/gotest/config_test.json")
	assert.Nil(t, err, "got unexpected error: %s", err)

	New(conf)

	assert.NotEqual(t, conf.Uchiwa.User, "*****", "Uchiwa user in private config shouldn't be masked")
	assert.NotEqual(t, conf.Uchiwa.Pass, "*****", "Uchiwa pass in private config shouldn't be masked")
	for i := range conf.Sensu {
		assert.NotEqual(t, conf.Sensu[i].User, "*****", "Sensu APIs user in private config shouldn't be masked")
		assert.NotEqual(t, conf.Sensu[i].Pass, "*****", "Sensu APIs pass in private config shouldn't be masked")
	}

	assert.Equal(t, PublicConfig.Uchiwa.User, "*****", "Uchiwa user in public config should be masked")
	assert.Equal(t, PublicConfig.Uchiwa.Pass, "*****", "Uchiwa pass in public config should be masked")
	for i := range PublicConfig.Sensu {
		assert.Equal(t, PublicConfig.Sensu[i].User, "*****", "Sensu APIs user in public config should be masked")
		assert.Equal(t, PublicConfig.Sensu[i].Pass, "*****", "Sensu APIs pass in public config should be masked")
	}

}
