// +build unit

package chanregistry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalChanRegistry(), "Global should have been set")

	var c *ChanRegistry
	SetGlobalChanRegistry(c)
	assert.Nil(t, GlobalChanRegistry(), "Global should be reset to nil")
}
