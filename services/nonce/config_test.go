// +build unit

package nonce

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Type(f)

	expected := redisOpt
	assert.Equal(t, expected, viper.GetString(typeViperKey), "Default")

	expected = "test"
	_ = os.Setenv(typeEnv, expected)
	assert.Equal(t, expected, viper.GetString(typeViperKey), "From Environment Variable")
	_ = os.Unsetenv(typeEnv)

	args := []string{
		"--nonce-manager-type=" + inMemoryOpt,
	}
	err := f.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = inMemoryOpt
	assert.Equal(t, expected, viper.GetString(typeViperKey), "From flag")
}
