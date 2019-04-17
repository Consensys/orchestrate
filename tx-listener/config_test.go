package listener

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

	os.Setenv("LISTENER_BLOCK_BACKOFF", "30s")
	expected = 30 * time.Second
	assert.Equal(t, expected, viper.GetDuration(name), "From Environment Variable")
	os.Unsetenv("LISTENER_BLOCK_BACKOFF")

	args := []string{
		"--listener-block-backoff=36s",
	}
	flgs.Parse(args)
	expected = 36 * time.Second
	assert.Equal(t, expected, viper.GetDuration(name), "From Flag")
}

func TestBlockLimit(t *testing.T) {
	name := "listener.block.limit"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BlockLimit(flgs)
	expected := int64(40)
	assert.Equal(t, expected, viper.GetInt64(name), "Default")

	os.Setenv("LISTENER_BLOCK_LIMIT", "45")
	expected = int64(45)
	assert.Equal(t, expected, viper.GetInt64(name), "From Environment Variable")
	os.Unsetenv("LISTENER_BLOCK_LIMIT")

	args := []string{
		"--listener-block-limit=60",
	}
	flgs.Parse(args)
	expected = int64(60)
	assert.Equal(t, expected, viper.GetInt64(name), "From Flag")
}

func TestTrackerDepth(t *testing.T) {
	name := "listener.tracker.depth"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	TrackerDepth(flgs)
	expected := int64(0)
	assert.Equal(t, expected, viper.GetInt64(name), "Default")

	os.Setenv("LISTENER_TRACKER_DEPTH", "45")
	expected = int64(45)
	assert.Equal(t, expected, viper.GetInt64(name), "From Environment Variable")
	os.Unsetenv("LISTENER_TRACKER_DEPTH")

	args := []string{
		"--listener-tracker-depth=60",
	}
	flgs.Parse(args)
	expected = int64(60)
	assert.Equal(t, expected, viper.GetInt64(name), "From Flag")
}
