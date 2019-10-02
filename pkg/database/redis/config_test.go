package redis

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	name := "redis.address"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	Address(flgs)
	expected := "localhost:6379"
	if viper.GetString(name) != expected {
		t.Errorf("RedisAddress #1: expected %q but got %q", expected, viper.GetString(name))
	}

	_ = os.Setenv("REDIS_ADDRESS", "127.0.0.1:6378")
	expected = "127.0.0.1:6378"
	if viper.GetString(name) != expected {
		t.Errorf("RedisAddress #2: expected %q but got %q", expected, viper.GetString(name))
	}

	args := []string{
		"--redis-address=127.0.0.1:6379",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = "127.0.0.1:6379"
	if viper.GetString(name) != expected {
		t.Errorf("RedisAddress #3: expected %q but got %q", expected, viper.GetString(name))
	}
}

func TestLockTimeout(t *testing.T) {
	name := "redis.lock.timeout"
	flgs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	LockTimeout(flgs)
	expected := 1500
	if viper.GetInt(name) != expected {
		t.Errorf("RedisLockTimeout #1: expected %v but got %v", expected, viper.GetInt(name))
	}

	_ = os.Setenv("REDIS_LOCKTIMEOUT", "2000")
	expected = 2000
	if viper.GetInt(name) != expected {
		t.Errorf("RedisLockTimeout #2: expected %v but got %v", expected, viper.GetInt(name))
	}

	args := []string{
		"--redis-lock-timeout=3000",
	}
	err := flgs.Parse(args)
	assert.NoError(t, err, "No error expected")

	expected = 3000
	if viper.GetInt(name) != expected {
		t.Errorf("RedisLockTimeout #3: expected %v but got %v", expected, viper.GetInt(name))
	}
}
