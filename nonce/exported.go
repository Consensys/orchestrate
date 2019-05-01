package nonce

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/nonce.git/nonce/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/nonce.git/nonce/redis"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	nc      Nonce
	initOnce = &sync.Once{}
)

// Init initializes Nonce
func Init() {
	initOnce.Do(func() {
		if nc != nil {
			return
		}

		switch viper.GetString(typeViperKey) {
		case "redis":
			// Initialize Redis Nonce Manager
			redis.Init()

			// Set Nonce
			nc = redis.GlobalNonce()
		case "mock":
			// Initialize Mock Nonce
			mock.Init()

			// Set Nonce
			nc = mock.GlobalNonce()
		default:
			log.WithFields(log.Fields{
				"type": viper.GetString(typeViperKey),
			}).Fatalf("nonce: unknown type")
		}
	})
}

// GlobalNonce returns global Sarama Nonce
func GlobalNonce() Nonce {
	return nc
}

// SetGlobalNonce sets global Sarama Nonce
func SetGlobalNonce(nonce Nonce) {
	nc = nonce
	log.Debug("nonce: set")
}
