package controllers

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/controllers/cooldown"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/controllers/creditor"
	maxbalance "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/controllers/max-balance"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
)

var (
	ctrl     ControlFunc
	initOnce = &sync.Once{}
)

// Init global controller
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if ctrl != nil {
			return
		}

		utils.InParallel(
			func() { faucet.Init(ctx) },
			func() { creditor.Init(ctx) },
			func() { maxbalance.Init(ctx) },
			func() { cooldown.Init(ctx) },
		)

		// Combine controls
		ctrl = CombineControls(
			creditor.Control,
			cooldown.Control,
			maxbalance.Control,
		)

		// Update global faucet
		faucet.SetGlobalFaucet(NewControlledFaucet(faucet.GlobalFaucet(), ctrl))
	})
}
