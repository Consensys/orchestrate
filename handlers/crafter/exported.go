package crafter

import (
	"context"
	"sync"

	"github.com/spf13/viper"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient/rpc"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/abi"
	registryclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/client"
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

		ec := ethclient.GlobalClient()

		// Create Handler
		handler = Crafter(registryclient.GlobalClient(), &abi.BaseCrafter{}, ec)

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
