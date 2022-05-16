package ristretto

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	_ = viper.BindEnv(CacheTTLViperKey, cacheTTLEnv)
	viper.SetDefault(CacheTTLViperKey, cacheDefault)
}

var (
	cacheTTLFlag     = "api-cache-ttl"
	CacheTTLViperKey = "api.cache.ttl"
	cacheDefault     = 0 * time.Second
	cacheTTLEnv      = "API_CACHE_TTL"
)

func Flags(f *pflag.FlagSet) {
	cacheDesc := fmt.Sprintf(`Cache TTL duration (Disabled by default). Environment variable: %q`, cacheTTLEnv)
	f.Duration(cacheTTLFlag, cacheDefault, cacheDesc)
	_ = viper.BindPFlag(CacheTTLViperKey, f.Lookup(cacheTTLFlag))
}
