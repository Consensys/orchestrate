package amount

import (
	"context"
	"math/big"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
)

var (
	ctrl     *Controller
	config   *Config
	initOnce = &sync.Once{}
)

// Init initialize BlackList Controller
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if ctrl != nil {
			return
		}

		// Set config if not yet set
		if config == nil {
			InitConfig(ctx)
		}

		// Initialize controller
		ctrl = NewController(config)

		log.WithFields(log.Fields{
			"controller": "amount",
		}).Info("faucet: controller ready")
	})
}

// InitConfig initialize configuration
func InitConfig(ctx context.Context) {
	amount, ok := big.NewInt(0).SetString(viper.GetString(faucetAmountViperKey), 10)
	if !ok {
		log.Fatalf("amount: invalid amount %q", viper.GetString(faucetAmountViperKey))
	}

	config = &Config{
		Amount: amount,
	}
}

// SetGlobalConfig sets global configuration
func SetGlobalConfig(c *Config) {
	config = c
}

// GlobalConfig returns global configuration
func GlobalConfig() *Config {
	return config
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
