package enricher

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	crc "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/contract-registry/client"
)

var (
	component = "handler.listener.enricher"
	handler   engine.HandlerFunc
	initOnce  = &sync.Once{}
)

// Init initialize handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Ethereum Client
		ethclient.Init(ctx)

		// Initialize Contract Registry Client
		crc.Init(ctx)

		// Create Handler
		handler = Enricher(crc.GlobalContractRegistryClient(), ethclient.GlobalClient())

		log.Infof("enricher: handler ready")
	})
}

// SetGlobalHandler sets global handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
