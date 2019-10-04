package handlers

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/handlers/dispatcher"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry/client"
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
