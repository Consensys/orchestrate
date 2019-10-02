package txlistener

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient"
	txlconfig "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/tx-listener/handler/base"
	txlhandler "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/tx-listener/handler/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/tx-listener/listener"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/enricher"
	envelopeloader "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/envelope/loader"
	receiptloader "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/loader/receipt"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/opentracing"
	producer "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/producer/tx-listener"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	server "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/http/healthcheck"
)

var (
	app       *common.App
	startOnce = &sync.Once{}
)

func init() {
	// Create app
	app = common.NewApp()
}

func startServer(ctx context.Context) {
	// Initialize server
	server.Init(ctx)

	// Register Healthcheck
	server.Enhance(healthcheck.HealthCheck(app))

	// Start Listening
	_ = server.ListenAndServe()
}

type serviceName string

func initHandlers(ctx context.Context) {
	common.InParallel(
		// Initialize Jaeger
		func() {
			ctx = context.WithValue(ctx, serviceName("service-name"), viper.GetString("jaeger.service.name"))
			opentracing.Init(ctx)
		},
		// Initialize store
		func() {
			envelopeloader.Init(ctx)
		},
		// Initialize enricher
		func() {
			enricher.Init(ctx)
		},
		// Initialize producer
		func() {
			producer.Init(ctx)
		},
	)
}

func initComponents(ctx context.Context) {
	common.InParallel(
		func() {
			engine.Init(ctx)
		},
		func() {
			initHandlers(ctx)
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
	engine.Register(receiptloader.Loader)
	engine.Register(opentracing.GlobalHandler())
	engine.Register(envelopeloader.GlobalHandler())
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
		app.SetReady(true)

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
