package redis

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const component = "nonce.redis"

var (
	nc       *Nonce
	initOnce = &sync.Once{}
)

// Init initializes Nonce
func Init() {
	initOnce.Do(func() {
		if nc != nil {
			return
		}

		pool := NewPool(viper.GetString(addressViperKey))

		// Initialize Nonce
		nc = NewNonce(pool, viper.GetInt(lockTimeoutViperKey))

		log.WithFields(log.Fields{
			"type": "redis",
		}).Info("nonce: ready")
	})
}

// GlobalNonce returns global Sarama Nonce
func GlobalNonce() *Nonce {
	return nc
}

// SetGlobalNonce sets global Sarama Nonce
func SetGlobalNonce(nonce *Nonce) {
	nc = nonce
}
