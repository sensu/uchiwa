package uchiwa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindDcFromString(t *testing.T) {
	dc := "foo"
	_, err := findDcFromString(&dc)
	assert.NotNil(t, err, "should return an error when a datacenter cannot be found")

	dc = "qux"
	_, err = findDcFromString(&dc)
	assert.Nil(t, err, "got unexpected error: %s", err)
}
