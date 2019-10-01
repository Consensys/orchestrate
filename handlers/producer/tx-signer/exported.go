package txsigner

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
)

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Gas Pricer Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Sync Producer
		broker.InitSyncProducer(ctx)

		// Create Handler
		handler = Producer(broker.GlobalSyncProducer())

		log.Infof("producer: handler ready")
	})
}

// SetGlobalHandler sets global Gas Pricer Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Gas Pricer handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
