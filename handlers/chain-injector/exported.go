package chaininjector

import (
	"context"
	"sync"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
)

const component = "handler.chaininjector"

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
	mux      = &sync.Mutex{}
)

// Init initialize Crafter Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Initialize eth client
		ethclient.Init(ctx)

		// Initialize chain-registry client
		registry.Init(ctx)

		multitenancyEnabled := viper.GetBool(multitenancy.EnabledViperKey)

		// Create Handler
		handler = engine.CombineHandlers(
			ChainInjector(multitenancyEnabled, registry.GlobalClient(), viper.GetString(registry.ChainRegistryURLViperKey)),
			ChainIDInjector(ethclient.GlobalClient()),
		)

		log.Infof("%s: handler ready", component)
	})
}

// SetGlobalHandler sets global Faucet Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	mux.Lock()
	handler = h
	mux.Unlock()
}

// GlobalHandler returns global Faucet handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
