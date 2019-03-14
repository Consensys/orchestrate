package worker

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSlots(t *testing.T) {
	name := "worker.slots"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Slots(flgs)
	expected := 20
	assert.Equal(t, expected, viper.GetInt(name), "Default")

	os.Setenv("WORKER_SLOTS", "125")
	expected = 125
	assert.Equal(t, expected, viper.GetInt(name), "From Environment Variable")
	os.Unsetenv("WORKER_SLOTS")

	args := []string{
		"--worker-slots=150",
	}
	flgs.Parse(args)
	expected = 150
	assert.Equal(t, expected, viper.GetInt(name), "From Flag")
}

func TestPartitions(t *testing.T) {
	name := "worker.partitions"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Partitions(flgs)
	expected := 50
	assert.Equal(t, expected, viper.GetInt(name), "Default")

	os.Setenv("WORKER_PARTITIONS", "125")
	expected = 125
	assert.Equal(t, expected, viper.GetInt(name), "From Environment Variable")
	os.Unsetenv("WORKER_PARTITIONS")

	args := []string{
		"--worker-partitions=150",
	}
	flgs.Parse(args)
	expected = 150
	assert.Equal(t, expected, viper.GetInt(name), "From Flag")
}
