package maxbalance

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
)

const component = "faucet.controllers.maxbalance"

var (
	ctrl     *Controller
	initOnce = sync.Once{}
)

// Init initialize BlackList Controller
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if ctrl != nil {
			return
		}

		// Initialize global MultiEthClient
		ethclient.Init(ctx)

		// Initialize controller
		ctrl = NewController(ethclient.GlobalClient())

		log.WithFields(log.Fields{
			"controller": "max-balance",
		}).Info("faucet: controller ready")
	})
}

// GlobalController returns global blacklist controller
func GlobalController() *Controller {
	return ctrl
}

// SetGlobalController sets global blacklist controller
func SetGlobalController(controller *Controller) {
	ctrl = controller
}

// Control allows to control a CreditFunc with global MaxBalance
func Control(f faucet.CreditFunc) faucet.CreditFunc {
	return ctrl.Control(f)
}
