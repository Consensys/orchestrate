package redis

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	name := "redis.url"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	URL(flgs)
	expected := "localhost:6379"
	if viper.GetString(name) != expected {
		t.Errorf("RedisURL #1: expected %q but got %q", expected, viper.GetString(name))
	}

	_ = os.Setenv("REDIS_URL", "127.0.0.1:6378")
	expected = "127.0.0.1:6378"
	if viper.GetString(name) != expected {
		t.Errorf("RedisURL #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--redis-url=127.0.0.1:6379",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "127.0.0.1:6379"
	if viper.GetString(name) != expected {
		t.Errorf("RedisURL #3: expected %q but got %q", expected, viper.GetString(name))
	}
}
