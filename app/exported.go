package app

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"

	log "github.com/sirupsen/logrus"

	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/logger"
	server "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
	txlconfig "gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/handler/base"
	txlhandler "gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/handler/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/listener"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/handlers/loader"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/handlers/producer"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/handlers/store"
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

func initComponents(ctx context.Context) {
	common.InParallel(
		func() {
			engine.Init(ctx)
		},
		func() {
			handlers.Init(ctx)
		},
		func() {
			broker.InitSyncProducer(ctx)
		},
		func() {
			listener.Init(ctx)
		},
	)
}

func registerHandlers() {
	// Generic handlers on every worker
	engine.Register(logger.Logger)

	// Specific handlers to tx-listener
	engine.Register(producer.GlobalHandler())
	engine.Register(loader.Loader)
	engine.Register(store.GlobalHandler())
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

		// Initialize all components of the server
		initComponents(cancelCtx)

		// Register all Handlers
		registerHandlers()

		// Indicate that application is ready
		// TODO: we need to update so ready can append when Consume has finished to Setup
		app.ready.Store(true)

		// Create handler
		conf, err := txlconfig.NewConfig()
		if err != nil {
			log.WithError(err).Fatalf("listener: could not load config")
		}
		h := txlhandler.NewHandler(engine.GlobalEngine(), broker.GlobalClient(), broker.GlobalSyncProducer(), conf)

		// Start Listening
		chains := ethclient.GlobalClient().Networks(cancelCtx)
		err = listener.Listen(cancelCtx, chains, h)
		if err != nil {
			log.WithError(err).Error("exiting loop with error")
		}
	})
}
