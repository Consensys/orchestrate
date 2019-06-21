package generator

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/keystore"
)

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Gas Estimator Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		common.InParallel(
			// Initialize keystore
			func() { keystore.Init(ctx) },
			// Initialize Sync Producer
			func() { broker.InitSyncProducer(ctx) },
		)

		// Create Handler
		handler = engine.CombineHandlers(
			WalletGenerator(keystore.GlobalKeyStore()),
			Producer(broker.GlobalSyncProducer()),
		)

		log.Infof("signer: handler ready")
	})
}

// SetGlobalHandler sets global Gas Estimator Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Gas Estimator handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
