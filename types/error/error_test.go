// +build unit

package error

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := New(1234, "Test Error").SetComponent("test-component")

	assert.Equal(t, "004D2@test-component: Test Error", err.Error(), "Error string should be correct")
	assert.Equal(t, uint64(1234), err.GetCode(), "Error code should be valid")
	assert.Equal(t, "test-component", err.GetComponent(), "Component should be valid")

	_ = err.SetMessage("%v %v", 3, fmt.Errorf("message test"))
	assert.Equal(t, "004D2@test-component: 3 message test", err.Error(), "Error string should be correct")
}

func TestExtendComponent(t *testing.T) {
	e := New(0, "test").ExtendComponent("foo")
	assert.Equal(t, "foo", e.GetComponent(), "Should set component correctly")

	e = e.ExtendComponent("bar")
	assert.Equal(t, "bar.foo", e.GetComponent(), "Should extend component correctly")
}
