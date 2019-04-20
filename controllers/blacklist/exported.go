package blacklist

import (
	"context"
	"math/big"
	"sync"

	ethcommon "github.com/ethereum/go-ethereum/common"
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

		// BlackList accounts by reading viper configuration
		blacklist := viper.GetStringSlice(faucetBlacklistViperKey)
		for _, bl := range blacklist {
			chainID, addr, err := utils.FromChainAccountKey(bl)
			if err != nil {
				logger.WithError(err).Fatalf("faucet: could not initialize controller")
			}
			ctrl.BlackList(chainID, addr)
		}
		logger.Info("faucet: controller ready")
	})
}

// BlackList  an account on a given chain
func BlackList(chainID *big.Int, address ethcommon.Address) {
	ctrl.BlackList(chainID, address)
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

// Control allows to control a CreditFunc with global Blacklist
func Control(f faucet.CreditFunc) faucet.CreditFunc {
	return ctrl.Control(f)
}
