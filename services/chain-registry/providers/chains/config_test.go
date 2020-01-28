package chains

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestProviderRefreshInterval(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	ProviderRefreshInterval(flgs)

	expected := time.Second
	assert.Equal(t, expected, viper.GetDuration(ProviderRefreshIntervalViperKey), "Default")

	_ = os.Setenv("CHAIN_REGISTRY_PROVIDER_CHAINS_REFRESH_INTERVAL", "30s")
	expected = 30 * time.Second
	assert.Equal(t, expected, viper.GetDuration(ProviderRefreshIntervalViperKey), "From Environment Variable")
	_ = os.Unsetenv("CHAIN_REGISTRY_PROVIDER_CHAINS_REFRESH_INTERVAL")

	args := []string{
		"--chain-registry-provider-chains-refresh-interval=36s",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = 36 * time.Second
	assert.Equal(t, expected, viper.GetDuration(ProviderRefreshIntervalViperKey), "From Flag")
}

func TestFlags(t *testing.T) {
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(flags)
	TestProviderRefreshInterval(t)
}
