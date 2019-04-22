package registry

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	viper.Set("abis", "ERC20[v0.1.3]:[{\"constant\":true,\"inputs\":[],\"name\":\"myFunction\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]")
	Init()

	m, err := GlobalRegistry().GetMethodByID("myFunction@ERC20[v0.1.3]")
	assert.NotNil(t, m, "Method should be available")
	assert.Nil(t, err, "Should not error")
}
