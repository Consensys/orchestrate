package redis

import (
	"sync"
	"time"

	healthz "github.com/heptiolabs/healthcheck"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const component = "nonce.redis"

var (
	nm       *NonceManager
	initOnce = &sync.Once{}
	checker  healthz.Check
)

// Init initializes Nonce
func Init() {
	initOnce.Do(func() {
		if nm != nil {
			return
		}

		redisURL := viper.GetString(URLViperKey)
		pool := NewPool(redisURL)

		// Initialize Nonce
		nm = NewNonceManager(pool, NewConfig())

		checker = healthz.TCPDialCheck(redisURL, time.Second*2)

		log.WithFields(log.Fields{
			"type": "redis",
		}).Info("nonce: ready")
	})
}

// GlobalNonceManager returns global NonceManager
func GlobalNonceManager() *NonceManager {
	return nm
}

func GlobalChecker() healthz.Check {
	return checker
}

// SetGlobalNonce sets global NonceManager
func SetGlobalNonceManager(m *NonceManager) {
	nm = m
}
