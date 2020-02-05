package handlers

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/handlers/dispatcher"
)

// Init handlers
func Init(ctx context.Context) {
	common.InParallel(
		func() {
			dispatcher.Init(ctx)
		},
	)
}
