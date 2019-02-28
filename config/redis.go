package config

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	redisAddressFlag     = "redis-address"
	redisAddressViperKey = "redis.address"
	redisAddressDefault  = "localhost:6379"
	redisAddressEnv      = "REDIS_ADDRESS"
)

// RedisAddress register a flag for Redis server address
func RedisAddress(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Address of Redis server to connect to.
Environment variable: %q`, redisAddressEnv)
	f.String(redisAddressFlag, redisAddressDefault, desc)
	viper.BindPFlag(redisAddressViperKey, f.Lookup(redisAddressFlag))
	viper.BindEnv(redisAddressViperKey, redisAddressEnv)
}

var (
	redisLockTimeoutFlag     = "redis-lock-timeout"
	redisLockTimeoutViperKey = "redis.lock.timeout"
	redisLockTimeoutDefault  = 1500
	redisLockTimeoutEnv      = "REDIS_LOCKTIMEOUT"
)

// RedisLockTimeout register a flag for Redis lock timeout
func RedisLockTimeout(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Redis lock timeout.
Environment variable: %q`, redisLockTimeoutEnv)
	f.Int(redisLockTimeoutFlag, redisLockTimeoutDefault, desc)
	viper.BindPFlag(redisLockTimeoutViperKey, f.Lookup(redisLockTimeoutFlag))
	viper.BindEnv(redisLockTimeoutViperKey, redisLockTimeoutEnv)
}

var (
	redisNonceExpirationTimeFlag     = "redis-nonce-expiration-date"
	redisNonceExpirationTimeViperKey = "redis.nonce.expiration.date"
	redisNonceExpirationTimeDefault  = 3
	redisNonceExpirationTimeEnv      = "REDIS_NONCE_EXPIRATION_DATE"
)

// RedisNonceExpirationTime register a flag for Redis nonce expiration time
func RedisNonceExpirationTime(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Redis nonce expiration time (duration in s).
Environment variable: %q`, redisNonceExpirationTimeEnv)
	f.Int(redisNonceExpirationTimeFlag, redisNonceExpirationTimeDefault, desc)
	viper.BindPFlag(redisNonceExpirationTimeViperKey, f.Lookup(redisNonceExpirationTimeFlag))
	viper.BindEnv(redisNonceExpirationTimeViperKey, redisNonceExpirationTimeEnv)
}
