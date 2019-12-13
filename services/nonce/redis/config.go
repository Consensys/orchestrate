package redis

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(URLViperKey, urlDefault)
	_ = viper.BindEnv(URLViperKey, urlEnv)
	viper.SetDefault(ExpirationViperKey, expirationDefault)
	_ = viper.BindEnv(ExpirationViperKey, expirationEnv)
}

// Register Redis flags
func Flags(f *pflag.FlagSet) {
	URL(f)
	Expiration(f)
}

const (
	urlFlag     = "redis-url"
	URLViperKey = "redis.url"
	urlDefault  = "localhost:6379"
	urlEnv      = "REDIS_URL"
)

// URL register a flag for Redis server URL
func URL(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`URL (address) of Redis server to connect to.
Environment variable: %q`, urlEnv)
	f.String(urlFlag, urlDefault, desc)
	_ = viper.BindPFlag(URLViperKey, f.Lookup(urlFlag))
}

const (
	expirationFlag     = "redis-expiration"
	ExpirationViperKey = "redis.expiration"
	expirationDefault  = 2 * time.Minute
	expirationEnv      = "REDIS_EXPIRATION"
)

// Expiration register a flag for Redis expiration
func Expiration(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Expiration for redis Key.
Environment variable: %q`, expirationEnv)
	f.Duration(expirationFlag, expirationDefault, desc)
	_ = viper.BindPFlag(ExpirationViperKey, f.Lookup(expirationFlag))
}

type Configuration struct {
	Expiration int
}

func NewConfig() *Configuration {
	return &Configuration{
		Expiration: int(viper.GetDuration(ExpirationViperKey).Milliseconds()),
	}
}
