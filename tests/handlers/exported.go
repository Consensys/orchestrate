package handlers

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/handlers/dispatcher"
)

// Init handlers
func Init(ctx context.Context) {
	utils.InParallel(
		func() {
			dispatcher.Init(ctx)
		},
	)
}
