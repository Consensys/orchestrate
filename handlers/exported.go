package handlers

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers/producer"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers/vault"
)

// Init inialize handlers
func Init(ctx context.Context) {
	common.InParallel(
		// Initialize Vault
		func() { vault.Init(ctx) },
		// Initialize Producer
		func() { producer.Init(ctx) },
	)
}
