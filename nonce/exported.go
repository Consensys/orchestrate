package nonce

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/service/nonce.git/nonce/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/nonce.git/nonce/redis"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	nc       Nonce
	initOnce = &sync.Once{}
)

// Init initializes Nonce
func Init(ctx context.Context) {
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
			mock.Init(ctx)

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
