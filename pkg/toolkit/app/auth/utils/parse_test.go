// +build unit

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAuth(t *testing.T) {
	auth, ok := ParseBearerToken("Bearer abdZ1==")
	assert.True(t, ok, "#1 ParseAuth should be valid")
	assert.Equal(t, "abdZ1==", auth, "#1 ParseAuth should parse correctly")

	auth, ok = ParseBearerToken("Basic abdZ1==")
	assert.False(t, ok, "#2 ParseAuth should not be valid")
	assert.Equal(t, "", auth, "#2 ParseAuth should parse correctly")
}
