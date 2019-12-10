package store

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRegistryType(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Type(flgs)

	expected := postgresOpt
	assert.Equal(t, expected, viper.GetString(TypeViperKey), "Default")

	expected = memoryOpt
	_ = os.Setenv(typeEnv, expected)
	assert.Equal(t, expected, viper.GetString(TypeViperKey), "From Environment Variable")
	_ = os.Unsetenv(typeEnv)

	args := []string{
		"--chain-registry-type=in-memory",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = memoryOpt
	assert.Equal(t, expected, viper.GetString(TypeViperKey), "From flag")
}
