package txlistener

import (
	"context"
	"fmt"
	"sync"

	"github.com/spf13/viper"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/metrics"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/key"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
	txlistener "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener"
)

var (
	app       = common.NewApp()
	startOnce = &sync.Once{}
)

// Start starts application
func Start(ctx context.Context) {
	startOnce.Do(func() {
		apiKey := viper.GetString(authkey.APIKeyViperKey)
		if apiKey != "" {
			// Inject authorization header in context for later authentication
			ctx = authutils.WithAuthorization(ctx, fmt.Sprintf("APIKey %v", apiKey))
		}

		cancelCtx, cancel := context.WithCancel(ctx)
		go metrics.StartServer(ctx, cancel, app.IsAlive, app.IsReady)

		// Initialize Tx-Listener
		txlistener.Init(cancelCtx)

		// Indicate that application is ready
		// TODO: we need to update so SetReady can be called when Consume has finished to Setup
		app.SetReady(true)

		// Start TxListener
		txlistener.Start(cancelCtx)
	})
}
