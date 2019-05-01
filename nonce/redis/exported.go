package redis

import (
	"github.com/spf13/viper"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	nc      *Nonce
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
			"type": "mock",
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
