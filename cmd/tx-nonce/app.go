package main

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/loader/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/logger"
	nonceattributor "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/nonce/attributor"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/offset"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/opentracing"
	producer "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/producer/tx-nonce"
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

	// Set Kafka Group value
	viper.Set("kafka.group", "group-nonce")
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
		// Initialize Jaeger tracer
		func() {
			ctx = context.WithValue(ctx, serviceName("service-name"), viper.GetString("jaeger.service.name"))
			opentracing.Init(ctx)
		},
		// Initialize Nonce manager
		func() {
			nonceattributor.Init(ctx)
		},
		// Initialize PrepareMsg
		func() {
			producer.Init(ctx)
		},
	)
}

func initComponents(ctx context.Context) {
	common.InParallel(
		// Initialize Engine
		func() {
			engine.Init(ctx)
		},
		// Initialize Handlers
		func() {
			initHandlers(ctx)
		},
		// Initialize ConsumerGroup
		func() {
			broker.InitConsumerGroup(ctx)
		},
	)
}

func registerHandlers() {
	// Generic handlers on every engine
	engine.Register(logger.Logger)
	engine.Register(sarama.Loader)
	engine.Register(offset.Marker)
	engine.Register(producer.GlobalHandler())
	engine.Register(opentracing.GlobalHandler())

	// Specific handlers tk Tx-Nonce worker
	engine.Register(nonceattributor.GlobalHandler())
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

		// Start consuming on topic tx-nonce
		err := broker.Consume(
			ctx,
			[]string{
				viper.GetString("kafka.topic.nonce"),
			},
			broker.NewEngineConsumerGroupHandler(engine.GlobalEngine()),
		)
		if err != nil {
			log.WithError(err).Error("worker: error on consumer")
		}
	})
}
