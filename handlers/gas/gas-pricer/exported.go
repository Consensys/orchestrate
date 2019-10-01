package gaspricer

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
)

var (
	component = "handler.gas-pricer"
	handler   engine.HandlerFunc
	initOnce  = &sync.Once{}
)

// Init initialize Gas Pricer Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Ethereum client
		ethclient.Init(ctx)

		// Create Handler
		handler = Pricer(ethclient.GlobalClient())

		log.Infof("gas-pricer: handler ready")
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
