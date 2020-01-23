package rpc

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMaxElapsedTime(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	MaxElapsedTime(f)
	expected := time.Hour
	assert.Equal(t, expected, viper.GetDuration(maxElapsedTimeViperKey), "Default")

	_ = os.Setenv(maxElapsedTimeEnv, "2s")
	expected, _ = time.ParseDuration("2s")
	assert.Equal(t, expected, viper.GetDuration(maxElapsedTimeViperKey), "From Environment Variable")
	_ = os.Unsetenv(maxElapsedTimeEnv)

	args := []string{
		"--ethclient-retry-maxelapsedtime=3s",
	}
	err := f.Parse(args)
	assert.Nil(t, err, "Parse Chain Registry flags should not error")
	expected, _ = time.ParseDuration("3s")
	assert.Equal(t, expected, viper.GetDuration(maxElapsedTimeViperKey), "From Flag")
}

func TestNewRetryConfig(t *testing.T) {
	r := NewRetryConfig()
	assert.NotNil(t, r, "Should get a retry config")
}

func TestNewConfig(t *testing.T) {
	r := NewConfig()
	assert.NotNil(t, r, "Should get a config")
}
