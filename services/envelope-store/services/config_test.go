package services

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestStoreType(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	StoreType(flgs)

	expected := "pg"
	assert.Equal(t, expected, viper.GetString(typeViperKey), "Default")

	expected = "env-store"
	_ = os.Setenv(typeEnv, expected)
	assert.Equal(t, expected, viper.GetString(typeViperKey), "From Environment Variable")
	_ = os.Unsetenv(typeEnv)

	args := []string{
		"--envelope-store=flag-store",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "flag-store"
	assert.Equal(t, expected, viper.GetString(typeViperKey), "From flag")
}
