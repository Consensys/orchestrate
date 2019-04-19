package redis

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(addressViperKey, addressDefault)
	viper.BindEnv(addressViperKey, addressEnv)
	viper.SetDefault(lockTimeoutViperKey, lockTimeoutDefault)
	viper.BindEnv(lockTimeoutViperKey, lockTimeoutEnv)
}

var (
	addressFlag     = "redis-address"
	addressViperKey = "redis.address"
	addressDefault  = "localhost:6379"
	addressEnv      = "REDIS_ADDRESS"
)

// Address register a flag for Redis server address
func Address(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Address of Redis server to connect to.
Environment variable: %q`, addressEnv)
	f.String(addressFlag, addressDefault, desc)
	viper.BindPFlag(addressViperKey, f.Lookup(addressFlag))
}

var (
	lockTimeoutFlag     = "redis-lock-timeout"
	lockTimeoutViperKey = "redis.lock.timeout"
	lockTimeoutDefault  = 1500
	lockTimeoutEnv      = "REDIS_LOCKTIMEOUT"
)

// LockTimeout register a flag for Redis lock timeout
func LockTimeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Redis lock timeout.
Environment variable: %q`, lockTimeoutEnv)
	f.Int(lockTimeoutFlag, lockTimeoutDefault, desc)
	viper.BindPFlag(lockTimeoutViperKey, f.Lookup(lockTimeoutFlag))
}
