package faucet

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet/sarama"
)

var (
	fct      Faucet
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if fct != nil {
			return
		}

		// Initialize Sarama Faucet
		sarama.Init(ctx)

		// Set Faucet
		fct = sarama.GlobalFaucet()
	})
}

// GlobalFaucet returns global Sarama Faucet
func GlobalFaucet() Faucet {
	return fct
}

// SetGlobalFaucet sets global Sarama Faucet
func SetGlobalFaucet(faucet Faucet) {
	fct = faucet
	log.Debug("faucet: set")
}
