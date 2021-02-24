package handlers

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/utils"
	"github.com/ConsenSys/orchestrate/tests/handlers/dispatcher"
)

// Init handlers
func Init(ctx context.Context) {
	utils.InParallel(
		func() {
			dispatcher.Init(ctx)
		},
	)
}
