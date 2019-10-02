package redis

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(addressViperKey, addressDefault)
	_ = viper.BindEnv(addressViperKey, addressEnv)
	viper.SetDefault(lockTimeoutViperKey, lockTimeoutDefault)
	_ = viper.BindEnv(lockTimeoutViperKey, lockTimeoutEnv)
	viper.SetDefault(expirationTimeViperKey, expirationTimeDefault)
	_ = viper.BindEnv(expirationTimeViperKey, expirationTimeEnv)
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
	_ = viper.BindPFlag(addressViperKey, f.Lookup(addressFlag))
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
	_ = viper.BindPFlag(lockTimeoutViperKey, f.Lookup(lockTimeoutFlag))
}

var (
	expirationTimeFlag     = "redis-nonce-expiration-time"
	expirationTimeViperKey = "redis.nonce.expiration.time"
	expirationTimeDefault  = 3
	expirationTimeEnv      = "REDIS_NONCE_EXPIRATION_TIME"
)

// ExpirationTime register expiration time flag
func ExpirationTime(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Redis nonce expiration time (duration in s).
Environment variable: %q`, expirationTimeEnv)
	f.Int(expirationTimeFlag, expirationTimeDefault, desc)
	_ = viper.BindPFlag(expirationTimeViperKey, f.Lookup(expirationTimeFlag))
}
