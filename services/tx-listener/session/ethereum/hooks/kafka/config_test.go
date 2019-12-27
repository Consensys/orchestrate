package kafka

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRegistryType(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	disableExternalTx(flgs)

	expected := false
	assert.Equal(t, expected, viper.GetBool(DisableExternalTxViperKey), "Default")

	_ = os.Setenv(disableExternalTxEnv, "true")
	assert.Equal(t, true, viper.GetBool(DisableExternalTxViperKey), "From Environment Variable")
	_ = os.Unsetenv(disableExternalTxEnv)

	args := []string{
		"--disable-external-tx=false",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = false
	assert.Equal(t, expected, viper.GetBool(DisableExternalTxViperKey), "From flag")
}
