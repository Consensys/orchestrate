package sender

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/envelope/storer"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	storeclient "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/envelope-store/client"
)

var (
	component = "handler.sender"
	handler   engine.HandlerFunc
	initOnce  = &sync.Once{}
)

// Init initialize Sender Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Context store
		storeclient.Init(ctx)

		// Initialize Ethereum client
		ethclient.Init(ctx)

		// Create Handler
		handler = engine.CombineHandlers(
			// Idempotency gate
			storer.TxAlreadySent(ethclient.GlobalClient(), storeclient.GlobalEnvelopeStoreClient()),
			// Sender
			Sender(ethclient.GlobalClient(), storeclient.GlobalEnvelopeStoreClient()),
		)

		log.Infof("sender: handler ready")
	})
}

// SetGlobalHandler sets global Sender Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Sender handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
