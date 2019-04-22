package engine

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSlots(t *testing.T) {
	name := "engine.slots"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Slots(flgs)
	expected := 20
	assert.Equal(t, expected, viper.GetInt(name), "Default")

	os.Setenv("WORKER_SLOTS", "125")
	expected = 125
	assert.Equal(t, expected, viper.GetInt(name), "From Environment Variable")
	os.Unsetenv("WORKER_SLOTS")

	args := []string{
		"--engine-slots=150",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "No error expected")

	expected = 150
	assert.Equal(t, expected, viper.GetInt(name), "From Flag")
}

func TestPartitions(t *testing.T) {
	name := "engine.partitions"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Partitions(flgs)
	expected := 50
	assert.Equal(t, expected, viper.GetInt(name), "Default")

	os.Setenv("WORKER_PARTITIONS", "125")
	expected = 125
	assert.Equal(t, expected, viper.GetInt(name), "From Environment Variable")
	os.Unsetenv("WORKER_PARTITIONS")

	args := []string{
		"--engine-partitions=150",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "No error expected")

	expected = 150
	assert.Equal(t, expected, viper.GetInt(name), "From Flag")
}

func TestTimeout(t *testing.T) {
	name := "engine.timeout"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Timeout(flgs)
	expected := 60 * time.Second
	assert.Equal(t, expected, viper.GetDuration(name), "Default")

	os.Setenv("WORKER_TIMEOUT", "12s")
	expected = 12 * time.Second
	assert.Equal(t, expected, viper.GetDuration(name), "From Environment Variable")
	os.Unsetenv("WORKER_TIMEOUT")

	args := []string{
		"--engine-timeout=100ms",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err, "No error expected")

	expected = 100 * time.Millisecond
	assert.Equal(t, expected, viper.GetDuration(name), "From Flag")
}
