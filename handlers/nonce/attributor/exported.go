package nonceattributor

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/nonce"
)

const component = "handler.nonce"

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
		nonce.Init(ctx)

		// Initialize the eth client
		ethclient.Init(ctx)

		// Create Nonce handler
		handler = Nonce(nonce.GlobalManager(), ethclient.GlobalClient())

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
