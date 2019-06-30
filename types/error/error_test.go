package error

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := New("Test Error").
		SetCode(1234).
		SetComponent("test-component")

	assert.Equal(t, "004D2@test-component: Test Error", err.Error(), "Error string should be correct")
	assert.Equal(t, uint64(1234), err.GetCode(), "Error code should be valid")
	assert.Equal(t, "test-component", err.GetComponent(), "Component should be valid")
}

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

func TestExtendComponent(t *testing.T) {
	e := New("test").ExtendComponent("foo")
	assert.Equal(t, "foo", e.GetComponent(), "Should set component correctly")

	e = e.ExtendComponent("bar")
	assert.Equal(t, "bar.foo", e.GetComponent(), "Should extend component correctly")
}
