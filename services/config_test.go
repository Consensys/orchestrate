package services

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestStoreType(t *testing.T) {
	name := "envelope-store.type"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	StoreType(flgs)

	expected := "pg"
	assert.Equal(t, expected, viper.GetString(name), "Default")

	os.Setenv("ENVELOPE_STORE", "env-store")
	expected = "env-store"
	assert.Equal(t, expected, viper.GetString(name), "From Environment Variable")
	os.Unsetenv("ENVELOPE_STORE")

	args := []string{
		"--envelope-store=flag-store",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "flag-store"
	assert.Equal(t, expected, viper.GetString(name), "From flag")
}
