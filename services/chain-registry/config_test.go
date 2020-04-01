// +build unit

package chainregistry

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

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
}
