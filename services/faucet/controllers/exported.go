package controllers

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/controllers/amount"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/controllers/blacklist"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/controllers/cooldown"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/controllers/creditor"
	maxbalance "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/controllers/max-balance"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/faucet/faucet"
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

		common.InParallel(
			func() { faucet.Init(ctx) },
			func() { maxbalance.Init(ctx) },
			func() { blacklist.Init(ctx) },
			func() { cooldown.Init(ctx) },
			func() { amount.Init(ctx) },
			func() { creditor.Init(ctx) },
		)

		// Combine controls
		ctrl = CombineControls(
			creditor.Control,
			blacklist.Control,
			cooldown.Control,
			amount.Control,
			maxbalance.Control,
		)

		// Update global faucet
		faucet.SetGlobalFaucet(NewControlledFaucet(faucet.GlobalFaucet(), ctrl))
	})
}

// SetControl sets global controller
func SetControl(control ControlFunc) {
	// Initialize Faucet
	faucet.Init(context.Background())
	ctrl = control

	// Update global faucet
	faucet.SetGlobalFaucet(NewControlledFaucet(faucet.GlobalFaucet(), ctrl))
}

// Control controls a credit function with global controller
func Control(f faucet.Faucet) faucet.Faucet {
	return NewControlledFaucet(f, ctrl)
}
