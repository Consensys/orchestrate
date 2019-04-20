package creditor

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/faucet.git/faucet"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
)

var (
	ctrl     *Controller
	initOnce = &sync.Once{}
)

// Init initialize BlackList Controller
func Init(ctx context.Context) {
	initOnce.Do(func() {
		// Initialize controller
		ctrl = NewController()

		// Enrich logger
		logger := log.WithFields(log.Fields{
			"controller": "blacklist",
		})

		// Set creditors
		for _, creditor := range viper.GetStringSlice(creditorAddressViperKey) {
			chainID, addr, err := utils.FromChainAccountKey(creditor)
			if err != nil {
				logger.WithError(err).Fatalf("faucet: could not initialize controller")
			}
			ctrl.SetCreditor(chainID, addr)
		}

		log.Info("faucet: controller ready")
	})
}

// GlobalController returns global blacklist controller
func GlobalController() *Controller {
	return ctrl
}

// SetGlobalController sets global blacklist controller
func SetGlobalController(controller *Controller) {
	initOnce.Do(func() {
		ctrl = controller
	})
}

// Control allows to control a CreditFunc with global CoolDown
func Control(f faucet.CreditFunc) faucet.CreditFunc {
	return ctrl.Control(f)
}
