package txlistener

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/metrics"
	txlistener "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener"
)

var (
	app       = common.NewApp()
	startOnce = &sync.Once{}
)

// Start starts application
func Start(ctx context.Context) {
	startOnce.Do(func() {
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
