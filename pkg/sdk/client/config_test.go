// +build unit

package client

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAPITarget(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(flgs)
	expected := urlDefault
	assert.Equal(t, expected, viper.GetString(URLViperKey), "Default")

	_ = os.Setenv(urlEnv, "env-api")
	expected = "env-api"
	assert.Equal(t, expected, viper.GetString(URLViperKey), "From Environment Variable")
	_ = os.Unsetenv(urlEnv)

	args := []string{
		"--api-url=flag-api",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "Parse API flags should not error")
	expected = "flag-api"
	assert.Equal(t, expected, viper.GetString(URLViperKey), "From Flag")
}

func TestFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
	assert.Equal(t, urlDefault, viper.GetString(URLViperKey), "Default")
}
