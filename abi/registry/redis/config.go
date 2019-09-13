package redis

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	redisURIFlag     = "contract-registry-redis-uri"
	redisURIViperKey = "contract.registry.redis.uri"
	redisURIDefault  = "localhost:6379"
	redisURIEnv      = "CONTRACT_REGISTRY_REDIS_URI"
	redisURIDesc     = fmt.Sprintf(`URI of the Redis contract registry to connect to.
	Environment variable: %q`, redisURIEnv)
)

var (
	redisMaxIdleConnFlag     = "contract-registry-redis-"
	redisMaxIdleConnViperKey = "contract.registry.redis.max.idle.conns"
	redisMaxIdleConnDefault  = 10000
	redisMaxIdleConnEnv      = "CONTRACT_REGISTRY_REDIS_MAX_IDLE_CONNS"
	redisMaxIdleConnDesc     = fmt.Sprintf(`Maximum number of idle connection in the redis pool.
	Environment variable: %q`, redisMaxIdleConnEnv)
)

var (
	redisMaxActiveConnFlag     = "contract-registry-redis-max-active-conns"
	redisMaxActiveConnViperKey = "contract.registry.redis.max.active.conns"
	redisMaxActiveConnDefault  = 20000
	redisMaxActiveConnEnv      = "CONTRACT_REGISTRY_REDIS_MAX_ACTIVE_CONNS"
	redisMaxActiveConnDesc     = fmt.Sprintf(`Maximum number of active connection  in the redis pool
	Environment variable: %q`, redisMaxActiveConnEnv)
)

var (
	redisMaxConnLifetimeFlag     = "contract-registry-redis-max-conn-lifetime"
	redisMaxConnLifetimeViperKey = "contract.registry.redis.max.conn.lifetime"
	redisMaxConnLifetimeDefault  = time.Duration(480) * time.Second
	redisMaxConnLifetimeEnv      = "CONTRACT_REGISTRY_REDIS_MAX_CONN_LIFETIME"
	redisMaxConnLifetimeDesc     = fmt.Sprintf(`Max lifetime of a redis connection  in the pool
	Environment variable: %q`, redisMaxConnLifetimeEnv)
)

var (
	redisIdleTimeoutFlag     = "contract-registry-redis-idle-timeout"
	redisIdleTimeoutViperKey = "contract.registry.redis.idle.timeout"
	redisIdleTimeoutDefault  = time.Duration(240) * time.Second
	redisIdleTimeoutEnv      = "CONTRACT_REGISTRY_REDIS_IDLE_TIMEOUT"
	redisIdleTimeoutDesc     = fmt.Sprintf(`Close connection after remaining idle for this duration
	Environment variable: %q`, redisIdleTimeoutEnv)
)

var (
	redisWaitFlag     = "contract-registry-redis-wait"
	redisWaitViperKey = "contract.registry.redis.wait"
	redisWaitDefault  = true
	redisWaitEnv      = "CONTRACT_REGISTRY_REDIS_WAIT"
	redisWaitDesc     = fmt.Sprintf(`If Wait is true and the pool is at the MaxActive limit, then Get() waits for a connection to be returned to the pool before returning.
	Environment variable: %q`, redisWaitEnv)
)

func init() {
	_ = viper.BindEnv(redisURIViperKey, redisURIEnv)
	viper.SetDefault(redisURIViperKey, redisURIDefault)

	_ = viper.BindEnv(redisMaxIdleConnViperKey, redisMaxIdleConnEnv)
	viper.SetDefault(redisMaxIdleConnViperKey, redisMaxIdleConnDefault)

	_ = viper.BindEnv(redisMaxActiveConnViperKey, redisMaxActiveConnEnv)
	viper.SetDefault(redisMaxActiveConnViperKey, redisMaxActiveConnDefault)

	_ = viper.BindEnv(redisMaxConnLifetimeViperKey, redisMaxConnLifetimeEnv)
	viper.SetDefault(redisMaxConnLifetimeViperKey, redisMaxConnLifetimeDefault)

	_ = viper.BindEnv(redisIdleTimeoutViperKey, redisIdleTimeoutEnv)
	viper.SetDefault(redisIdleTimeoutViperKey, redisIdleTimeoutDefault)

	_ = viper.BindEnv(redisWaitViperKey, redisWaitEnv)
	viper.SetDefault(redisWaitViperKey, redisWaitDefault)
}

// InitFlags batch registers all flags
func InitFlags(f *pflag.FlagSet) {
	f.String(redisURIFlag, redisURIDefault, redisURIDesc)
	f.Int(redisMaxIdleConnFlag, redisMaxIdleConnDefault, redisMaxIdleConnDesc)
	f.Int(redisMaxActiveConnFlag, redisMaxActiveConnDefault, redisMaxActiveConnDesc)
	f.Duration(redisMaxConnLifetimeFlag, redisMaxConnLifetimeDefault, redisMaxConnLifetimeDesc)
	f.Duration(redisIdleTimeoutFlag, redisIdleTimeoutDefault, redisIdleTimeoutDesc)
	f.Bool(redisWaitFlag, redisWaitDefault, redisWaitDesc)
}

// Config returns a PoolConfig object set from viper keys
func Config() *PoolConfig {
	return &PoolConfig{
		URI:             viper.GetString(redisURIViperKey),
		MaxIdle:         viper.GetInt(redisMaxIdleConnViperKey),
		MaxActive:       viper.GetInt(redisMaxActiveConnViperKey),
		MaxConnLifetime: viper.GetDuration(redisMaxConnLifetimeViperKey),
		IdleTimeout:     viper.GetDuration(redisIdleTimeoutViperKey),
		Wait:            viper.GetBool(redisWaitViperKey),
	}
}
