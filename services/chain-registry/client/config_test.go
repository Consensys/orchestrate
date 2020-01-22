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
	ChainRegistryURL(flgs)
	expected := "localhost:8081"
	assert.Equal(t, expected, viper.GetString(ChainRegistryURLViperKey), "Default")

	_ = os.Setenv(chainRegistryURLEnv, "env-chain-registry")
	expected = "env-chain-registry"
	assert.Equal(t, expected, viper.GetString(ChainRegistryURLViperKey), "From Environment Variable")
	_ = os.Unsetenv(chainRegistryURLEnv)

	args := []string{
		"--chain-registry-url=flag-chain-registry",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "Parse Chain Registry flags should not error")
	expected = "flag-chain-registry"
	assert.Equal(t, expected, viper.GetString(ChainRegistryURLViperKey), "From Flag")
}

func TestChainProxyURL(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	ChainProxyURL(f)
	expected := "localhost:8080"
	assert.Equal(t, expected, viper.GetString(ChainProxyURLViperKey), "Default")

	_ = os.Setenv(chainProxyURLEnv, "env-chain-proxy")
	expected = "env-chain-proxy"
	assert.Equal(t, expected, viper.GetString(ChainProxyURLViperKey), "From Environment Variable")
	_ = os.Unsetenv(chainProxyURLEnv)

	args := []string{
		"--chain-proxy-url=flag-chain-proxy",
	}
	err := f.Parse(args)
	assert.Nil(t, err, "Parse Chain Proxy flags should not error")
	expected = "flag-chain-proxy"
	assert.Equal(t, expected, viper.GetString(ChainProxyURLViperKey), "From Flag")
}

func TestFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
	assert.Equal(t, chainProxyURLDefault, viper.GetString(ChainProxyURLViperKey), "Default")
	assert.Equal(t, chainRegistryURLDefault, viper.GetString(ChainRegistryURLViperKey), "Default")
}
