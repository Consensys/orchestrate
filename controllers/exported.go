package controllers

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/controllers/amount"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/controllers/blacklist"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/controllers/cooldown"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/controllers/creditor"
	maxbalance "gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/controllers/max-balance"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/faucet"
)

var (
	ctrl     ControlFunc
	initOnce = &sync.Once{}
)

// Init intiliaze global controller
func Init(ctx context.Context) {
	initOnce.Do(func() {
		// Initialize Faucet
		faucet.Init(ctx)

		// Initialize Controllers
		wg := &sync.WaitGroup{}
		wg.Add(5)
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

		log.Info("faucet: ready")
	})
}

// SetControl sets global controller
func SetControl(control ControlFunc) {
	initOnce.Do(func() {
		ctrl = control
	})
}

// Control controls a credit function with global controller
func Control(f faucet.Faucet) faucet.Faucet {
	return NewControlledFaucet(f, ctrl)
}
