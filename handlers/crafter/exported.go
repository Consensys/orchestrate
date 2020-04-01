package crafter

import (
	"context"
	"sync"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	registryclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry/client"
)

const component = "handler.crafter"

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Crafter Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize Registry Client
		registryclient.Init(ctx, viper.GetString(registryclient.ContractRegistryURLViperKey))

		// Create Handler
		handler = Crafter(registryclient.GlobalClient(), &abi.BaseCrafter{})

		log.Infof("crafter: handler ready")
	})
}

// SetGlobalHandler sets global Faucet Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Faucet handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
