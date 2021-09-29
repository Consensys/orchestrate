package handlers

import (
	"context"

	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/tests/handlers/dispatcher"
)

// Init handlers
func Init(ctx context.Context) {
	utils.InParallel(
		func() {
			dispatcher.Init(ctx)
		},
	)
}
