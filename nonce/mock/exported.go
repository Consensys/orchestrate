package mock

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	nc       *Nonce
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if nc != nil {
			return
		}

		// Initialize Faucet
		nc = NewNonce()

		log.WithFields(log.Fields{
			"type": "mock",
		}).Info("nonce: ready")
	})
}

// GlobalFaucet returns global Mock Nonce
func GlobalNonce() *Nonce {
	return nc
}

// SetGlobalFaucet sets global Mock Nonce
func SetGlobalNonce(nonce *Nonce) {
	nc = nonce
}
