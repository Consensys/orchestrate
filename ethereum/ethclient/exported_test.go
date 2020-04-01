// +build unit

package ethclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	Init(context.Background())
	assert.NotNil(t, GlobalClient(), "Global should have been set")

	var c Client
	SetGlobalClient(c)
	assert.Nil(t, GlobalClient(), "Global should be reset to nil")
}
