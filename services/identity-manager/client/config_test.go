// +build unit

package client

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestKeyManagerTarget(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(flgs)
	expected := urlDefault
	assert.Equal(t, expected, viper.GetString(URLViperKey), "Default")

	_ = os.Setenv(urlEnv, "env-identity-manager")
	expected = "env-identity-manager"
	assert.Equal(t, expected, viper.GetString(URLViperKey), "From Environment Variable")
	_ = os.Unsetenv(urlEnv)

	args := []string{
		"--identity-manager-url=flag-identity-manager",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "Parse Identity Manager flags should not error")
	expected = "flag-identity-manager"
	assert.Equal(t, expected, viper.GetString(URLViperKey), "From Flag")
}

func TestFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
	assert.Equal(t, urlDefault, viper.GetString(URLViperKey), "Default")
}
