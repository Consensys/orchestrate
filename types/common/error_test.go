package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := &Error{
		Type:    uint64(18),
		Message: "Test Error",
	}

	assert.Equal(t, "Error #18: Test Error", err.Error(), "Error message should be valid")
}
