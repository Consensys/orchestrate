package faucet

import (
	"context"
	"sync"

	"github.com/spf13/viper"
	chaininjector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/chain-injector"
	handlerfaucet "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/faucet"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/controllers"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/faucet/faucet"
)

const component = "handler.account-generator.faucet"

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

		// Initialize Controlled Faucet
		controllers.Init(ctx)

		// Initialize chain-registry client
		registry.Init(ctx)

		// Create Handler
		handler = engine.CombineHandlers(
			chaininjector.ChainUUIDHandlerWithoutAbort(registry.GlobalClient(), viper.GetString(registry.ChainRegistryURLViperKey)),
			handlerfaucet.Faucet(faucet.GlobalFaucet(), registry.GlobalClient()),
		)

		log.Infof("%s: handler ready", component)
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
