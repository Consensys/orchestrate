// +build unit

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	assert.True(t, *Bool(true))
	assert.False(t, *Bool(false))
}
