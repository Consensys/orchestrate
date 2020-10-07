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

		cfg := NewConfig(viper.GetViper())
		pool, err := NewPool(cfg)
		if err != nil {
			log.WithError(err).Fatalf("could not connect to redis server")
		}

		// Initialize Nonce
		nm = NewNonceManager(pool, cfg)

		checker = healthz.Timeout(nm.Ping, time.Second*2)

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
