package app

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	server "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http/healthcheck"
)

var (
	app       *App
	startOnce = &sync.Once{}
)

func init() {
	// Create app
	app = NewApp()
}

func startServer(ctx context.Context) {
	// Initialize server
	server.Init(ctx)

	// Register Healthcheck
	server.Enhance(healthcheck.HealthCheck(app))

	// Start Listening
	_ = server.ListenAndServe()
}

// Start starts application
func Start(ctx context.Context) {
	startOnce.Do(func() {

		cancelCtx, cancel := context.WithCancel(ctx)
		go func() {
			// Start Server
			startServer(ctx)
			cancel()
		}()

		// Indicate that application is ready
		// TODO: we need to update so ready can append when Consume has finished to Setup
		app.ready.Store(true)

		// Code below is an example for the unique purpose of illustrating this boilerplate
		ticker := time.NewTicker(time.Second)
	runningLoop:
		for {
			select {
			case <-cancelCtx.Done():
				log.Info("Leaving Loop")
				break runningLoop
			case <-ticker.C:
				log.Info("Ticking")
			}
		}
	})
}
