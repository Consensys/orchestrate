package store

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitRegistry(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	InitRegistry(f)

	var expected []string
	assert.Equal(t, expected, viper.GetStringSlice(InitViperKey), "Default")

	_ = os.Setenv(initEnv, "test1 test2")
	assert.Equal(t, []string{"test1", "test2"}, viper.GetStringSlice(InitViperKey), "From Environment Variable")
	_ = os.Unsetenv(initEnv)
}

func TestFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
}
