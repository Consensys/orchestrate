package vault

import (
	"context"
	"sync"

	chaininjector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/chain-injector"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/signer"
	generator "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault/wallet-generator"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init handlers
func Init(ctx context.Context) {
	initOnce.Do(func() {
		common.InParallel(
			// Initialize keystore
			func() { signer.Init(ctx) },
			// Initialize Sync Producer
			func() { generator.Init(ctx) },
			// Initialize ChainID injector
			func() { chaininjector.Init(ctx) },
		)

		signerHandler := engine.CombineHandlers(
			chaininjector.GlobalHandler(),
			signer.GlobalHandler(),
		)

		handler = Vault(signerHandler, generator.GlobalHandler())
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
