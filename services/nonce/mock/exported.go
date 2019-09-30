package mock

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	nm       *NonceManager
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if nm != nil {
			return
		}

		// Initialize Faucet
		nm = NewNonceManager()

		log.WithFields(log.Fields{
			"type": "mock",
		}).Info("nonce: ready")
	})
}

// GlobalFaucet returns global Mock Nonce
func GlobalNonceManager() *NonceManager {
	return nm
}

// SetGlobalFaucet sets global Mock Nonce
func SetGlobalNonceManager(m *NonceManager) {
	nm = m
}
