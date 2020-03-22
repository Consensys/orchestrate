package handlers

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/handlers/dispatcher"
)

// Init handlers
func Init(ctx context.Context) {
	utils.InParallel(
		func() {
			dispatcher.Init(ctx)
		},
	)
}
