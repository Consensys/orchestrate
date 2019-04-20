package controllers

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/controllers/blacklist"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/controllers/cooldown"
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
		wg := &sync.WaitGroup{}
		wg.Add(3)
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
		wg.Wait()

		// Combine controls
		ctrl = CombineControls(blacklist.Control, cooldown.Control, maxbalance.Control)

		log.Info("faucet: controllers ready")
	})
}

// SetControl sets global controller
func SetControl(control ControlFunc) {
	ctrl = control
}

// Control controls a credit function with global controller
func Control(credit faucet.CreditFunc) faucet.CreditFunc {
	return ctrl(credit)
}
