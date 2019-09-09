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
	m        Manager
	initOnce = &sync.Once{}
)

// Init initializes Nonce
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if m != nil {
			return
		}

		switch viper.GetString(typeViperKey) {
		case "redis":
			// Initialize Redis Nonce Manager
			redis.Init()

			// Set Nonce
			m = redis.GlobalNonceManager()
		case "mock":
			// Initialize Mock Nonce
			mock.Init(ctx)

			// Set Nonce
			m = mock.GlobalNonceManager()
		default:
			log.WithFields(log.Fields{
				"type": viper.GetString(typeViperKey),
			}).Fatalf("nonce: unknown storage type")
		}
	})
}

// GlobalManager returns globalNonceManager
func GlobalManager() Manager {
	return m
}

// SetGlobalManager sets global Manager
func SetGlobalManager(mngr Manager) {
	m = mngr
}
