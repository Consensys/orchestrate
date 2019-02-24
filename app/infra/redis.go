package infra

import (
	"sync"

	"github.com/spf13/viper"
	infRedis "gitlab.com/ConsenSys/client/fr/core-stack/infra/redis.git"
)

func initRedis(infra *Infra, wait *sync.WaitGroup) {
	infra.NonceManager = infRedis.NewNonceManager(viper.GetString("redis.address"), viper.GetInt("redis.lock.timeout"))
	wait.Done()
}
