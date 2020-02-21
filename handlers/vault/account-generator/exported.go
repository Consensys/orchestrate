package generator

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/account-generator/account"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/account-generator/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
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
			func() { account.Init(ctx) },
			func() { faucet.Init(ctx) },
		)

		// Create Handler
		handler = engine.CombineHandlers(
			account.GlobalHandler(),
			faucet.GlobalHandler(),
		)

		log.Infof("account-generator: handler ready")
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
