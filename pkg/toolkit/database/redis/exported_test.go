// +build unit

package redis

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	viper.Reset()
	Init()
	assert.NotNil(t, GlobalClient(), "Faucet should have been set")

	var n *Client
	SetGlobalClient(n)
	assert.Nil(t, GlobalClient(), "Global should be reset to nil")
}
