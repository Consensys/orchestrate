package base

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
