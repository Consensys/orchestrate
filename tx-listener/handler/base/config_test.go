package base

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestStartDefault(t *testing.T) {
	name := "listener.start-default"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	StartDefault(flgs)

	expected := "oldest"
	assert.Equal(t, expected, viper.GetString(name), "Default")

	os.Setenv("LISTENER_START_DEFAULT", "latest")
	expected = "latest"
	assert.Equal(t, expected, viper.GetString(name), "Env")
	os.Unsetenv("LISTENER_START_DEFAULT")

	args := []string{
		"--listener-start-default=123",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "No error expected")

	expected = "123"
	assert.Equal(t, expected, viper.GetString(name), "Flag")
}

func TestStart(t *testing.T) {
	name := "listener.start"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Start(flgs)

	expected := []string{}
	assert.Equal(t, expected, viper.GetStringSlice(name), "Default")

	os.Setenv("LISTENER_START", "3:oldest")
	expected = []string{"3:oldest"}
	assert.Equal(t, expected, viper.GetStringSlice(name), "Env")
	os.Unsetenv("LISTENER_START")

	args := []string{
		"--listener-start=2:latest",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "No error expected")

	expected = []string{"2:latest"}
	assert.Equal(t, expected, viper.GetStringSlice(name), "Flag")
}

func TestConfig(t *testing.T) {
	viper.Set("listener.start", []string{"2:oldest", "3:128-3"})
	viper.Set("listener.start-default", "genesis")

	conf, err := NewConfig()
	assert.Nil(t, err, "No error expected")
	assert.Equal(t, int64(0), conf.Start.Default.BlockNumber, "Default start should be correct")
	assert.Equal(t, int64(-2), conf.Start.Positions["2"].BlockNumber, "BlockNumber on chain 2 should be correct")
	assert.Equal(t, int64(128), conf.Start.Positions["3"].BlockNumber, "BlockNumber on chain 3 should be correct")
}
