// +build unit

package http

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.NotNil(t, cfg)
}

func TestHostname(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	hostname(flgs)
	assert.Equal(t, hostnameDefault, viper.GetString(hostnameViperKey))

	expected := "localhost"
	_ = os.Setenv(hostnameEnv, expected)
	assert.Equal(t, expected, viper.GetString(hostnameViperKey))
	_ = os.Unsetenv(hostnameEnv)

	expected = "127.0.0.1"
	args := []string{
		fmt.Sprintf("--%s=%s", hostnameFlag, expected),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, expected, viper.GetString(hostnameViperKey))
}

func TestHTTPPort(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	port(flgs)
	assert.Equal(t, httpPortDefault, viper.GetInt(httpPortViperKey))

	expected := "9999"
	_ = os.Setenv(httpPortEnv, expected)
	assert.Equal(t, expected, viper.GetString(httpPortViperKey))
	_ = os.Unsetenv(httpPortEnv)

	expected = "9989"
	args := []string{
		fmt.Sprintf("--%s=%s", httpPortFlag, expected),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, expected, viper.GetString(httpPortViperKey))
}

func TestMetricsHostname(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	metricsHostname(flgs)
	assert.Equal(t, metricsHostnameDefault, viper.GetString(metricsHostnameViperKey))

	expected := "localhost"
	_ = os.Setenv(metricsHostnameEnv, expected)
	assert.Equal(t, expected, viper.GetString(metricsHostnameViperKey))
	_ = os.Unsetenv(metricsHostnameEnv)

	expected = "127.0.0.1"
	args := []string{
		fmt.Sprintf("--%s=%s", metricsHostnameFlag, expected),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, expected, viper.GetString(metricsHostnameViperKey))
}

func TestMetricsPort(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	metricsPort(flgs)
	assert.Equal(t, metricsPortDefault, viper.GetInt(metricsPortViperKey))

	expectedOne := "9999"
	_ = os.Setenv(metricsPortEnv, expectedOne)
	assert.Equal(t, expectedOne, viper.GetString(metricsPortViperKey))
	_ = os.Unsetenv(metricsPortEnv)

	expectedTwo := "9989"
	args := []string{
		fmt.Sprintf("--%s=%s", metricsPortFlag, expectedTwo),
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")
	assert.Equal(t, expectedTwo, viper.GetString(metricsPortViperKey))
}

func TestFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
}
