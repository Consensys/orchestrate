// +build unit

package key

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, checker, "auth Manager should have been set")

	var c *Key
	SetGlobalChecker(c)
	assert.Nil(t, GlobalChecker(), "Global Auth should be reset to nil")
}
