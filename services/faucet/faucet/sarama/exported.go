package sarama

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/broker/sarama"
)

var (
	component = "faucet.sarama"
	fct       *Faucet
	initOnce  = &sync.Once{}
)

// Init initializes Faucet
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if fct != nil {
			return
		}

		// Initialize Producer
		broker.InitSyncProducer(ctx)

		// Control Faucet
		fct = NewFaucet(broker.GlobalSyncProducer())

		log.WithFields(log.Fields{
			"type": "sarama",
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
