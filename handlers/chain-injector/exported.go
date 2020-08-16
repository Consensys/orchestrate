package chaininjector

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient/rpc"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
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

		// Create Handler
		handler = engine.CombineHandlers(
			ChainUUIDHandler(registry.GlobalClient(), viper.GetString(registry.ChainRegistryURLViperKey)),
			ChainIDInjectorHandler(ethclient.GlobalClient()),
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
