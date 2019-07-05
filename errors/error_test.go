package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorf(t *testing.T) {
	err := Errorf("Test %q", "msg")
	assert.Equal(t, "Test \"msg\"", err.GetMessage(), "Error message should be valid")
	assert.Equal(t, "FF000", err.Hex(), "Hex code should be correct")
}

func TestFromError(t *testing.T) {
	assert.Nil(t, FromError(nil), "From nil error should be nil")
	e := FromError(fmt.Errorf("test"))
	assert.Equal(t, "test", e.GetMessage(), "Error message should be correct")
	assert.Equal(t, "FF000", e.Hex(), "Hex code should be correct")

	e2 := FromError(e)
	assert.Equal(t, e, e2, "Should behave as flat pass on internal errors")
}

func TestIsErrorClass(t *testing.T) {
	assert.True(t, isErrorClass(271120, 270336), "A 42310 should be a 42000")
	assert.False(t, isErrorClass(270336, 270848), "A 42310 should not be a 42200")
	assert.False(t, isErrorClass(270336, 0), "A 42000 should not be a 00000")
	assert.False(t, isErrorClass(0, 270336), "A 00000 should not be a 42000")
	assert.False(t, isErrorClass(275216, 270336), "A 43310 should not be a 42000")
}
