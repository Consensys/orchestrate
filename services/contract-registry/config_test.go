// +build unit

package contractregistry

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)


func TestFlags(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Type(flgs)

	assert.Equal(t, abiDefault, viper.GetStringSlice(ABIViperKey), "")

	expected := []string{"<contract>:<abi>:<bytecode>:<deployedBytecode>"}
	_ = os.Setenv(abiEnv, "<contract>:<abi>:<bytecode>:<deployedBytecode>")
	assert.Equal(t, expected, viper.GetStringSlice(ABIViperKey), "From Environment Variable")
	_ = os.Unsetenv(abiEnv)

	args := []string{
		"--abi=<contract>:<abi>:<bytecode>:<deployedBytecode>",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = []string{"<contract>:<abi>:<bytecode>:<deployedBytecode>"}
	assert.Equal(t, expected, viper.GetStringSlice(ABIViperKey), "From flag")
}
