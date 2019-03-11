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

func TestRedisLockTimeout(t *testing.T) {
	name := "redis.lock.timeout"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RedisLockTimeout(flgs)
	expected := 1500
	if viper.GetInt(name) != expected {
		t.Errorf("RedisLockTimeout #1: expected %v but got %v", expected, viper.GetInt(name))
	}

	os.Setenv("REDIS_LOCKTIMEOUT", "2000")
	expected = 2000
	if viper.GetInt(name) != expected {
		t.Errorf("RedisLockTimeout #2: expected %v but got %v", expected, viper.GetInt(name))
	}

	args := []string{
		"--redis-lock-timeout=3000",
	}
	flgs.Parse(args)
	expected = 3000
	if viper.GetInt(name) != expected {
		t.Errorf("RedisLockTimeout #3: expected %v but got %v", expected, viper.GetInt(name))
	}
}
