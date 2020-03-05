package cooldown

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
)

const component = "faucet.controllers.cooldown"

var (
	ctrl     *Controller
	initOnce = &sync.Once{}
)

// Init initialize BlackList Controller
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if ctrl != nil {
			return
		}

		// Initialize controller
		ctrl = NewController()

		log.WithFields(log.Fields{
			"controller": "cooldown",
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

// Control allows to control a CreditFunc with global CoolDown
func Control(f faucet.CreditFunc) faucet.CreditFunc {
	return ctrl.Control(f)
}
