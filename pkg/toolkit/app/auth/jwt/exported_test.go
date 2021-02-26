// +build unit

package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	var c *JWT
	SetGlobalChecker(c)
	assert.Nil(t, GlobalChecker(), "Global Auth should be reset to nil")
}
