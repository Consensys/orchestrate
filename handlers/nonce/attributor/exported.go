package nonceattributor

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/engine"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/nonce"
)

const component = "handler.nonce.attributor"

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
		ec := ethclient.GlobalClient()
		handler = Nonce(nonce.GlobalManager(), ec)

		log.Infof("%s: handler ready", component)
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
