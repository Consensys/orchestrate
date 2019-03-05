package infra

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	infRedis "gitlab.com/ConsenSys/client/fr/core-stack/infra/redis.git"
)

var (
	redisNonceExpirationTimeFlag     = "redis-nonce-expiration-time"
	redisNonceExpirationTimeViperKey = "redis.nonce.expiration.time"
	redisNonceExpirationTimeDefault  = 3
	redisNonceExpirationTimeEnv      = "REDIS_NONCE_EXPIRATION_TIME"
)

// RedisNonceExpirationTime register a flag for Redis nonce expiration time
func RedisNonceExpirationTime(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Redis nonce expiration time (duration in s).
Environment variable: %q`, redisNonceExpirationTimeEnv)
	f.Int(redisNonceExpirationTimeFlag, redisNonceExpirationTimeDefault, desc)
	viper.BindPFlag(redisNonceExpirationTimeViperKey, f.Lookup(redisNonceExpirationTimeFlag))
	viper.BindEnv(redisNonceExpirationTimeViperKey, redisNonceExpirationTimeEnv)
}

func initNonce(infra *Infra, wait *sync.WaitGroup) {
	infra.NonceManager = infRedis.NewNonceManager(viper.GetString("redis.address"), viper.GetInt("redis.lock.timeout"))
	log.Infof("infra-nonce: nonce manager ready")
	wait.Done()
}
