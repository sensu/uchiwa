package uchiwa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClient(t *testing.T) {
	_, err := GetClient("foo", "qux")
	assert.Nil(t, err, "got unexpected error: %s", err)
}
