package tessera

import (
	"context"
	"sync"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/tessera"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
)

const component = "handler.tessera"

var (
	handler  engine.HandlerFunc
	initOnce = &sync.Once{}
)

// Init initialize Gas Estimator Handler
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		tessera.Init(ctx)

		// Create Handler
		handler = StoreRaw(
			tessera.GlobalClient(),
			viper.GetString(registry.ChainRegistryURLViperKey),
		)

		log.Infof("tessera: handler ready")
	})
}

// SetGlobalHandler sets global Gas Estimator Handler
func SetGlobalHandler(h engine.HandlerFunc) {
	handler = h
}

// GlobalHandler returns global Gas Estimator handler
func GlobalHandler() engine.HandlerFunc {
	return handler
}
