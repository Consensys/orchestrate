package registry

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestProviderRefreshInterval(t *testing.T) {
	name := "provider.refreshInterval"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	ProviderRefreshInterval(flgs)
	expected := time.Second
	assert.Equal(t, expected, viper.GetDuration(name), "Default")

	_ = os.Setenv("PROVIDER_REFRESHINTERVAL", "30s")
	expected = 30 * time.Second
	assert.Equal(t, expected, viper.GetDuration(name), "From Environment Variable")
	_ = os.Unsetenv("PROVIDER_REFRESHINTERVAL")

	args := []string{
		"--provider-refreshInterval=36s",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = 36 * time.Second
	assert.Equal(t, expected, viper.GetDuration(name), "From Flag")
}

func TestFlags(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(flags)
	TestProviderRefreshInterval(t)
}
