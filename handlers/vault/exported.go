package vault

import (
	"context"
	"sync"

	chaininjector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/chain-injector"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/signer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init handlers
func Init(ctx context.Context) {
	initOnce.Do(func() {
		utils.InParallel(
			// Initialize keystore
			func() { signer.Init(ctx) },
			// Initialize ChainID injector
			func() { chaininjector.Init(ctx) },
		)

		signerHandler := engine.CombineHandlers(
			chaininjector.GlobalHandler(),
			multitenancy.GlobalHandler(),
			signer.GlobalHandler(),
		)

		handler = Vault(signerHandler)
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
