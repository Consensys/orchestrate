package redis

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(URLViperKey, urlDefault)
	_ = viper.BindEnv(URLViperKey, urlEnv)
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
