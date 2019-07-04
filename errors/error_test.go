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
	assert.Equal(t, "FF000", e.Hex(), "Hex code should be correct")

	e2 := FromError(e)
	assert.Equal(t, e, e2, "Should behave as flat pass on internal errors")
}

func TestIsErrorClass(t *testing.T) {
	assert.True(t, isErrorClass(271120, dataErrCode), "Hex 42310 should be a data error")
	assert.False(t, isErrorClass(dataErrCode, solidityErrCode), "Data error should not be a solidity error")
	assert.False(t, isErrorClass(dataErrCode, 0), "Hex 00000 should not be a data error")
	assert.False(t, isErrorClass(0, dataErrCode), "Hex 00000 should not be a data error")
	assert.False(t, isErrorClass(275216, dataErrCode), "Hex 43310 should not be a data error")
}
