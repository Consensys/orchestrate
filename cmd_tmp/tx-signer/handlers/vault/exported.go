package vault

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers/vault/signer"
	generator "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers/vault/wallet-generator"
)

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init inialize handlers
func Init(ctx context.Context) {
	initOnce.Do(func() {
		common.InParallel(
			// Initialize keystore
			func() { signer.Init(ctx) },
			// Initialize Sync Producer
			func() { generator.Init(ctx) },
		)

		handler = Vault(signer.GlobalHandler(), generator.GlobalHandler())
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
