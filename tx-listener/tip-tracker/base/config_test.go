package base

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestTrackerDepth(t *testing.T) {
	name := "listener.tracker.depth"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Depth(flgs)
	expected := int64(0)
	assert.Equal(t, expected, viper.GetInt64(name), "Default")

	os.Setenv("LISTENER_TRACKER_DEPTH", "45")
	expected = int64(45)
	assert.Equal(t, expected, viper.GetInt64(name), "From Environment Variable")
	os.Unsetenv("LISTENER_TRACKER_DEPTH")

	args := []string{
		"--listener-tracker-depth=60",
	}
	err := flgs.Parse(args)
	assert.Nil(t, err)

	expected = int64(60)
	assert.Equal(t, expected, viper.GetInt64(name), "From Flag")
}
