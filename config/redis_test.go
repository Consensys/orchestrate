package config

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func TestRedisAddress(t *testing.T) {
	name := "redis.address"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RedisAddress(flgs)
	expected := "localhost:6379"
	if viper.GetString(name) != expected {
		t.Errorf("RedisAddress #1: expected %q but got %q", expected, viper.GetString(name))
	}

	os.Setenv("REDIS_ADDRESS", "127.0.0.1:6378")
	expected = "127.0.0.1:6378"
	if viper.GetString(name) != expected {
		t.Errorf("RedisAddress #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--redis-address=127.0.0.1:6379",
	}
	flgs.Parse(args)
	expected = "127.0.0.1:6379"
	if viper.GetString(name) != expected {
		t.Errorf("RedisAddress #3: expected %q but got %q", expected, viper.GetString(name))
	}
}
