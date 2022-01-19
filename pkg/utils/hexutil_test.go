package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHexString(t *testing.T) {
	assert.True(t, IsHexString("0x12"))
	assert.False(t, IsHexString("0xGG"))
}
