package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayIndexOf_exists(t *testing.T) {
	idx := ArrayIndexOf([]string{"a", "b", "c"}, "b")
	assert.Equal(t, 1, idx)
}

func TestArrayIndexOf_notExists(t *testing.T) {
	idx := ArrayIndexOf([]string{"a", "b", "c"}, "d")
	assert.Equal(t, -1, idx)
}

func TestArrayIndexOf_error(t *testing.T) {
	idx := ArrayIndexOf([]string{"a", "b", "c"}, "d")
	assert.Equal(t, -1, idx)
}

func TestArrayIntersection_matched(t *testing.T) {
	result := ArrayIntersection([]string{"a", "b", "c"}, []string{"b", "c", "e"})
	assert.Len(t, result, 2)
	assert.Equal(t, []string{"b", "c"}, result.([]string))
}

func TestArrayIndexOf_empty(t *testing.T) {
	result := ArrayIntersection([]string{"a", "b", "c"}, []string{"d", "e"})
	assert.Len(t, result, 0)
}
