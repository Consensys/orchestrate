// +build unit

package client

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestChainRegistryTarget(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(flgs)
	expected := chainRegistryURLDefault
	assert.Equal(t, expected, viper.GetString(ChainRegistryURLViperKey), "Default")

	_ = os.Setenv(ChainRegistryURLEnv, "env-chain-registry")
	expected = "env-chain-registry"
	assert.Equal(t, expected, viper.GetString(ChainRegistryURLViperKey), "From Environment Variable")
	_ = os.Unsetenv(ChainRegistryURLEnv)

	args := []string{
		"--chain-registry-url=flag-chain-registry",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "Parse Chain Registry flags should not error")
	expected = "flag-chain-registry"
	assert.Equal(t, expected, viper.GetString(ChainRegistryURLViperKey), "From Flag")
}

func TestFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
	assert.Equal(t, chainRegistryURLDefault, viper.GetString(ChainRegistryURLViperKey), "Default")
}
