package uchiwa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayIntersection(t *testing.T) {
	var array1 []string
	var array2 []string

	found := arrayIntersection(array1, array2)
	assert.Equal(t, false, found, "if both arrays are empty, it should return false")

	array1 = []string{"foo", "bar"}
	found = arrayIntersection(array1, array2)
	assert.Equal(t, false, found, "if one array is empty, it should return false")

	array2 = []string{"baz", "qux"}
	found = arrayIntersection(array1, array2)
	assert.Equal(t, false, found, "it should return false is none of the elements in the arrays are shared")

	array2 = append(array2, "foo")
	found = arrayIntersection(array1, array2)
	assert.Equal(t, true, found, "it should return true if at least one element is shared between the arrays")
}
