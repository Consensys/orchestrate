package faucet

import (
	"context"
	"sync"

	faucetscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet/scheduler"

	log "github.com/sirupsen/logrus"
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
		faucetscheduler.Init()

		// Set Faucet
		fct = faucetscheduler.GlobalFaucet()
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
