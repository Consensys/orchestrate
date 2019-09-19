package crafter

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	registryclient "gitlab.com/ConsenSys/client/fr/core-stack/service/contract-registry.git/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/abi/crafter"
)

var (
	component = "handler.crafter"
	handler   engine.HandlerFunc
	initOnce  = &sync.Once{}
)

// Init initialize Crafter Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Create crafter
		crafter.Init()

		// Initialize Registry Client
		registryclient.Init(ctx)

		// Create Handler
		handler = Crafter(registryclient.GlobalContractRegistryClient(), crafter.GlobalCrafter())

		log.Infof("crafter: handler ready")
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
