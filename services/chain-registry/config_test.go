package chainregistry

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestProxyAddress(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	ProxyAddress(f)
	expected := ":8080"
	assert.Equal(t, expected, viper.GetString(ProxyAddressViperKey), "Default")

	_ = os.Setenv(proxyAddressEnv, "proxy-address-env")
	expected = "proxy-address-env"
	assert.Equal(t, expected, viper.GetString(ProxyAddressViperKey), "From Environment Variable")
	_ = os.Unsetenv(proxyAddressEnv)

	args := []string{
		"--chain-proxy-addr=proxy-address-flag",
	}
	err := f.Parse(args)
	assert.Nil(t, err, "Parse Chain Proxy flags should not error")
	expected = "proxy-address-flag"
	assert.Equal(t, expected, viper.GetString(ProxyAddressViperKey), "From Flag")
}

func TestRegistryAddress(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RegistryAddress(f)
	expected := ":8081"
	assert.Equal(t, expected, viper.GetString(AddressViperKey), "Default")

	_ = os.Setenv(addressEnv, "chain-registry-addr-env")
	expected = "chain-registry-addr-env"
	assert.Equal(t, expected, viper.GetString(AddressViperKey), "From Environment Variable")
	_ = os.Unsetenv(addressEnv)

	args := []string{
		"--chain-registry-addr=chain-registry-addr-flag",
	}
	err := f.Parse(args)
	assert.Nil(t, err, "Parse Chain Registry flags should not error")
	expected = "chain-registry-addr-flag"
	assert.Equal(t, expected, viper.GetString(AddressViperKey), "From Flag")
}

func TestProvidersThrottleDuration(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	ProvidersThrottleDuration(f)
	expected := time.Second
	assert.Equal(t, expected, viper.GetDuration(ProvidersThrottleDurationViperKey), "Default")

	_ = os.Setenv(providersThrottleDurationEnv, "2s")
	expected, _ = time.ParseDuration("2s")
	assert.Equal(t, expected, viper.GetDuration(ProvidersThrottleDurationViperKey), "From Environment Variable")
	_ = os.Unsetenv(providersThrottleDurationEnv)

	args := []string{
		"--providers-throttle-duration=3s",
	}
	err := f.Parse(args)
	assert.Nil(t, err, "Parse Chain Registry flags should not error")
	expected, _ = time.ParseDuration("3s")
	assert.Equal(t, expected, viper.GetDuration(ProvidersThrottleDurationViperKey), "From Flag")
}

func TestNewConfig(t *testing.T) {
	c := NewConfig()
	assert.NotNil(t, c, "Should get config")
}

func TestFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
	TestProxyAddress(t)
	TestRegistryAddress(t)
	TestProvidersThrottleDuration(t)
}
