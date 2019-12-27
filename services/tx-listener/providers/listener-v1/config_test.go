package listenerv1

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestBlockBackoff(t *testing.T) {
	name := "listener.block.backoff"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BlockBackoff(flgs)
	expected := time.Second
	assert.Equal(t, expected, viper.GetDuration(name), "Default")

	_ = os.Setenv("LISTENER_BLOCK_BACKOFF", "30s")
	expected = 30 * time.Second
	assert.Equal(t, expected, viper.GetDuration(name), "From Environment Variable")
	_ = os.Unsetenv("LISTENER_BLOCK_BACKOFF")

	args := []string{
		"--listener-block-backoff=36s",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = 36 * time.Second
	assert.Equal(t, expected, viper.GetDuration(name), "From Flag")
}

func TestBlockLimit(t *testing.T) {
	name := "listener.block.limit"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BlockLimit(flgs)
	expected := int64(40)
	assert.Equal(t, expected, viper.GetInt64(name), "Default")

	_ = os.Setenv("LISTENER_BLOCK_LIMIT", "45")
	expected = int64(45)
	assert.Equal(t, expected, viper.GetInt64(name), "From Environment Variable")
	_ = os.Unsetenv("LISTENER_BLOCK_LIMIT")

	args := []string{
		"--listener-block-limit=60",
	}

	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = int64(60)
	assert.Equal(t, expected, viper.GetInt64(name), "From Flag")
}

func TestTrackerDepth(t *testing.T) {
	name := "listener.tracker.depth"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Depth(flgs)
	expected := int64(0)
	assert.Equal(t, expected, viper.GetInt64(name), "Default")

	_ = os.Setenv("LISTENER_TRACKER_DEPTH", "45")
	expected = int64(45)
	assert.Equal(t, expected, viper.GetInt64(name), "From Environment Variable")
	_ = os.Unsetenv("LISTENER_TRACKER_DEPTH")

	args := []string{
		"--listener-tracker-depth=60",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = int64(60)
	assert.Equal(t, expected, viper.GetInt64(name), "From Flag")
}

func TestStartDefault(t *testing.T) {
	name := "listener.start-default"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	StartDefault(flgs)

	expected := "oldest"
	assert.Equal(t, expected, viper.GetString(name), "Default")

	_ = os.Setenv("LISTENER_START_DEFAULT", "latest")
	expected = "latest"
	assert.Equal(t, expected, viper.GetString(name), "Env")
	_ = os.Unsetenv("LISTENER_START_DEFAULT")

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

	var expected []string
	assert.Equal(t, expected, viper.GetStringSlice(name), "Default")

	_ = os.Setenv("LISTENER_START", "3:oldest")
	expected = []string{"3:oldest"}
	assert.Equal(t, expected, viper.GetStringSlice(name), "Env")
	_ = os.Unsetenv("LISTENER_START")

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
