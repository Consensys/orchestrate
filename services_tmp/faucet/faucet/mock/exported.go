package mock

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	fct      *Faucet
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if fct != nil {
			return
		}

		// Initialize Faucet
		fct = NewFaucet()

		log.WithFields(log.Fields{
			"type": "mock",
		}).Info("faucet: ready")
	})
}

// GlobalFaucet returns global Sarama Faucet
func GlobalFaucet() *Faucet {
	return fct
}

// SetGlobalFaucet sets global Sarama Faucet
func SetGlobalFaucet(faucet *Faucet) {
	fct = faucet
}
