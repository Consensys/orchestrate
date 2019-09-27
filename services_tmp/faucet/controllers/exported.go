package controllers

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/controllers/amount"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/controllers/blacklist"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/controllers/cooldown"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/controllers/creditor"
	maxbalance "gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/controllers/max-balance"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/faucet.git/faucet"
)

var (
	ctrl     ControlFunc
	initOnce  = &sync.Once{}
)

// Init intiliaze global controller
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if ctrl != nil {
			return
		}

		// Initialize Controllers
		wg := &sync.WaitGroup{}
		wg.Add(6)

		// Initialize Faucet
		go func() {
			faucet.Init(ctx)
			wg.Done()
		}()

		// Initialize Controllers
		go func() {
			maxbalance.Init(ctx)
			wg.Done()
		}()
		go func() {
			blacklist.Init(ctx)
			wg.Done()
		}()
		go func() {
			cooldown.Init(ctx)
			wg.Done()
		}()
		go func() {
			amount.Init(ctx)
			wg.Done()
		}()
		go func() {
			creditor.Init(ctx)
			wg.Done()
		}()
		wg.Wait()

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
