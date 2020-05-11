// +build unit

package client

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestTransactionSchedulerTarget(t *testing.T) {
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(flgs)
	expected := txSchedulerURLDefault
	assert.Equal(t, expected, viper.GetString(txSchedulerURLViperKey), "Default")

	_ = os.Setenv(txSchedulerURLEnv, "env-transaction-scheduler")
	expected = "env-transaction-scheduler"
	assert.Equal(t, expected, viper.GetString(txSchedulerURLViperKey), "From Environment Variable")
	_ = os.Unsetenv(txSchedulerURLEnv)

	args := []string{
		"--transaction-scheduler-url=flag-transaction-scheduler",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "Parse Transaction Scheduler flags should not error")
	expected = "flag-transaction-scheduler"
	assert.Equal(t, expected, viper.GetString(txSchedulerURLViperKey), "From Flag")
}

func TestFlags(t *testing.T) {
	f := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Flags(f)
	assert.Equal(t, txSchedulerURLDefault, viper.GetString(txSchedulerURLViperKey), "Default")
}
