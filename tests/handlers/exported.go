package handlers

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/common"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/contract-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/tests/handlers/dispatcher"
)

// Init handlers
func Init(ctx context.Context) {
	common.InParallel(
		func() {
			dispatcher.Init(ctx)
		},
		// Initialize the registryClient
		func() {
			registry.Init(ctx)
		},
	)
}
