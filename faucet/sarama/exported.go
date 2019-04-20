package sarama

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
)

var (
	fct      *Faucet
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		// Initialize Producer
		broker.InitSyncProducer(ctx)

		// Control Faucet
		fct = NewFaucet(broker.SyncProducer())

		log.Info("faucet: ready")
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
