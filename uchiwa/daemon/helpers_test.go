package daemon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBoolFromInterface(t *testing.T) {
	i := map[string]interface{}{"foo": true}

	_, err := GetBoolFromInterface(i)
	assert.NotNil(t, err)

	b, err := GetBoolFromInterface(i["foo"])
	assert.Nil(t, err)
	assert.Equal(t, b, true)
}

func TestGetMapFromInterface(t *testing.T) {
	i := map[string]interface{}{"foo": "vodka"}
	m := GetMapFromInterface(i)
	assert.Equal(t, "vodka", m["foo"])
}

func TestStringInArray(t *testing.T) {
	var item string
	var array []string

	found := StringInArray(item, array)
	assert.Equal(t, false, found, "if item and array are both empty, it should return false")

	item = "foo"
	found = StringInArray(item, array)
	assert.Equal(t, false, found, "if array is empty, it should return false")

	array = []string{"bar", "qux"}
	found = StringInArray(item, array)
	assert.Equal(t, false, found, "it should return false if the item isn't found in the array")

	array = append(array, "foo")
	found = StringInArray(item, array)
	assert.Equal(t, true, found, "it should return true if the item is found in the array")
}
