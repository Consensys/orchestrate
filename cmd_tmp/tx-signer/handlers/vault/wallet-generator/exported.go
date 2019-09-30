package generator

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers/vault/wallet-generator/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers/vault/wallet-generator/wallet"
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
			func() { wallet.Init(ctx) },
			func() { faucet.Init(ctx) },
		)

		// Create Handler
		handler = engine.CombineHandlers(
			wallet.GlobalHandler(),
			faucet.GlobalHandler(),
		)

		log.Infof("wallet-generator: handler ready")
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
