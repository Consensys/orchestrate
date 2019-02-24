package infra

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	infRedis "gitlab.com/ConsenSys/client/fr/core-stack/infra/redis.git"
)

func initNonce(infra *Infra, wait *sync.WaitGroup) {
	infra.NonceManager = infRedis.NewNonceManager(viper.GetString("redis.address"), viper.GetInt("redis.lock.timeout"))
	log.Infof("infra-nonce: nonce manager ready")
	wait.Done()
}
