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

	_, err = LoadConfig("../test/gotest/config_test.json")
	assert.Nil(t, err, "got unexpected error: %s", err)
}
