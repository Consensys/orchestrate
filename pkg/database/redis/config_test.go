// +build unit

package redis

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRedisHost(t *testing.T) {
	name := "redis.host"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	URL(flgs)
	expected := "localhost"
	if viper.GetString(name) != expected {
		t.Errorf("RedisHOST #1: expected %q but got %q", expected, viper.GetString(name))
	}

	_ = os.Setenv("REDIS_HOST", "127.0.0.1")
	expected = "127.0.0.1"
	if viper.GetString(name) != expected {
		t.Errorf("RedisHOST #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--redis-host=127.0.0.1",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "127.0.0.1"
	if viper.GetString(name) != expected {
		t.Errorf("RedisHOST #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestRedisPort(t *testing.T) {
	name := "redis.port"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	URL(flgs)
	expected := "6379"
	if viper.GetString(name) != expected {
		t.Errorf("RedisPORT #1: expected %q but got %q", expected, viper.GetString(name))
	}

	_ = os.Setenv("REDIS_PORT", "6378")
	expected = "6378"
	if viper.GetString(name) != expected {
		t.Errorf("RedisPORT #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--redis-port=6377",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "6377"
	if viper.GetString(name) != expected {
		t.Errorf("RedisPORT #3: expected %q but got %q", expected, viper.GetString(name))
	}
}
