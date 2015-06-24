package uchiwa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceIntersection(t *testing.T) {
	var a1, a2 []string

	found := SliceIntersection(a1, a2)
	assert.Equal(t, false, found, "if both slices are empty, it should return false")

	a1 = []string{"foo", "bar"}
	found = SliceIntersection(a1, a2)
	assert.Equal(t, false, found, "if one slice is empty, it should return false")

	a2 = []string{"baz", "qux"}
	found = SliceIntersection(a1, a2)
	assert.Equal(t, false, found, "it should return false is none of the elements in the slices are shared")

	a2 = append(a2, "foo")
	found = SliceIntersection(a1, a2)
	assert.Equal(t, true, found, "it should return true if at least one element is shared between the slices")
}

func TestMergeStringSlice(t *testing.T) {
	var a1, a2 []string

	slice := MergeStringSlices(a1, a2)
	assert.Equal(t, []string(nil), slice, "if both slices are empty, it should return an empty slice")

	a1 = []string{"1", "2"}
	slice = MergeStringSlices(a1, a2)
	assert.Equal(t, a1, slice, "if one slice is empty, it should return the other slice")

	a2 = []string{"2"}
	slice = MergeStringSlices(a1, a2)
	assert.Equal(t, a1, slice, "if one slice is empty, it should return the other slice")

	a2 = []string{"2", "3"}
	slice = MergeStringSlices(a1, a2)
	assert.Equal(t, []string{"1", "2", "3"}, slice, "if one slice is empty, it should return the other slice")

}
