package redis

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault(addressViperKey, addressDefault)
	_ = viper.BindEnv(addressViperKey, addressEnv)
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
