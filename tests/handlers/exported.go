package handlers

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"
	registryClient "gitlab.com/ConsenSys/client/fr/core-stack/service/contract-registry.git/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/handlers/dispatcher"
)

// Init handlers
func Init(ctx context.Context) {
	common.InParallel(
		func() {
			dispatcher.Init(ctx)
		},
		// Initialize the registryClient
		func() {
			registryClient.Init(ctx)
		},
	)
}
