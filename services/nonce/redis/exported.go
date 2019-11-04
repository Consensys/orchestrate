package redis

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const component = "nonce.redis"

var (
	nm       *NonceManager
	initOnce = &sync.Once{}
)

// Init initializes Nonce
func Init() {
	initOnce.Do(func() {
		if nm != nil {
			return
		}

		pool := NewPool(viper.GetString(URLViperKey))

		// Initialize Nonce
		nm = NewNonceManager(pool)

		log.WithFields(log.Fields{
			"type": "redis",
		}).Info("nonce: ready")
	})
}

// GlobalNonceManager returns global NonceManager
func GlobalNonceManager() *NonceManager {
	return nm
}

// SetGlobalNonce sets global NonceManager
func SetGlobalNonceManager(m *NonceManager) {
	nm = m
}
