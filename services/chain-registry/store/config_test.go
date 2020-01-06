package store

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRegistryType(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Type(f)

	expected := postgresOpt
	assert.Equal(t, expected, viper.GetString(TypeViperKey), "Default")

	expected = memoryOpt
	_ = os.Setenv(typeEnv, expected)
	assert.Equal(t, expected, viper.GetString(TypeViperKey), "From Environment Variable")
	_ = os.Unsetenv(typeEnv)

	args := []string{
		"--chain-registry-type=in-memory",
	}
	err := f.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = memoryOpt
	assert.Equal(t, expected, viper.GetString(TypeViperKey), "From flag")
}

func TestInitRegistry(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	InitRegistry(f)

	expected := []string{}
	assert.Equal(t, expected, viper.GetStringSlice(InitViperKey), "Default")

	_ = os.Setenv(initEnv, "test1 test2")
	assert.Equal(t, []string{"test1", "test2"}, viper.GetStringSlice(InitViperKey), "From Environment Variable")
	_ = os.Unsetenv(initEnv)

	args := []string{
		"--chain-registry-init=test2",
		"--chain-registry-init=test3",
	}
	err := f.Parse(args)
	assert.NoError(t, err, "No error expected")

	assert.Equal(t, []string{"test2", "test3"}, viper.GetStringSlice(InitViperKey), "From flag")
}

func TestFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
	TestInitRegistry(t)
	TestRegistryType(t)
}
