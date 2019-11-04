package creditor

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
)

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

		// Enrich logger
		logger := log.WithFields(log.Fields{
			"controller": "creditor",
		})

		// Set creditors
		for _, creditor := range viper.GetStringSlice(creditorAddressViperKey) {
			chainID, addr, err := utils.FromChainAccountKey(creditor)
			if err != nil {
				logger.WithError(err).Fatalf("faucet: could not initialize controller")
			}
			ctrl.SetCreditor(chainID, addr)
		}

		logger.Info("faucet: controller ready")
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
