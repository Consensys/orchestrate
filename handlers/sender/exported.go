package sender

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	grpcStore "gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/store/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
)

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Sender Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Context store
		grpcStore.Init(ctx)

		// Initialize Ethereum client
		ethclient.Init(ctx)

		// Create Handler
		handler = Sender(ethclient.GlobalClient(), grpcStore.GlobalEnvelopeStore())

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
