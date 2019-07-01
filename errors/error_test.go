package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorf(t *testing.T) {
	err := Errorf("Test %q", "msg")
	assert.Equal(t, "Test \"msg\"", err.GetMessage(), "Error message should be valid")
}

func TestFromError(t *testing.T) {
	assert.Nil(t, FromError(nil), "From nil error should be nil")
	e := FromError(fmt.Errorf("test"))
	assert.Equal(t, "test", e.GetMessage(), "Error message should be correct")

	e2 := FromError(e)
	assert.Equal(t, e, e2, "Should behave as flat pass on internal errors")
}
