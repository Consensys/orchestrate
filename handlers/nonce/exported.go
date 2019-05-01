package nonce

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/nonce.git/nonce"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Nonce Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize the nonce manager
		nonce.Init()

		// Initialize the eth client
		ethclient.Init(ctx)

		// Create Handler
		handler = NonceHandler(nonce.GlobalNonce(), ethclient.GlobalMultiClient().PendingNonceAt)

		log.Infof("logger: handler ready")
	})
}

// SetGlobalHandler sets global Faucet Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Faucet handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
