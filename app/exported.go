package app

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/loader"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/offset"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-nonce.git/handlers/nonce"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-nonce.git/handlers/producer"
	"sync"

	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	server "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-nonce.git/handlers"
)

var (
	app       *App
	startOnce = &sync.Once{}
)

func init() {
	// Create app
	app = NewApp()

	// Set Kafka Group value
	viper.Set("kafka.group", "group-nonce")
}

func startServer(ctx context.Context) {
	// Initialize server
	server.Init(ctx)

	// Register Healthcheck
	server.Enhance(healthcheck.HealthCheck(app))

	// Start Listening
	server.ListenAndServe()
}

func initComponents(ctx context.Context) {
	wg := sync.WaitGroup{}

	// Initialize Engine
	wg.Add(1)
	go func() {
		engine.Init(ctx)
		wg.Done()
	}()

	// Initialize Handlers
	wg.Add(1)
	go func() {
		handlers.Init(ctx)
		wg.Done()
	}()

	// Initialize ConsumerGroup
	wg.Add(1)
	go func() {
		broker.InitConsumerGroup(ctx)
		wg.Done()
	}()

	// Wait for engine and handlers to be ready
	wg.Wait()
}

func registerHandlers() {
	wg := sync.WaitGroup{}

	// Register handlers on engine
	wg.Add(1)
	go func() {
		// Generic handlers on every engine
		engine.Register(logger.Logger)
		engine.Register(loader.Loader)
		engine.Register(offset.Marker)
		engine.Register(producer.GlobalHandler())

		// Specific handlers tk Tx-Nonce worker
		engine.Register(nonce.GlobalHandler())

		wg.Done()
	}()

	// Wait for ConsumerGroup & Engine to be ready
	wg.Wait()
}

// Start starts application
func Start(ctx context.Context) {
	startOnce.Do(func() {
		ctx, cancel := context.WithCancel(ctx)
		go func() {
			// Start Server
			startServer(ctx)
			cancel()
		}()

		// Initialize all components of the server
		initComponents(ctx)

		// Register all Handlers
		registerHandlers()

		// Indicate that application is ready
		// TODO: we need to update so ready can append when Consume has finished to Setup
		app.ready.Store(true)

		// Start consuming on topic tx-nonce
		broker.Consume(
			ctx,
			[]string{
				viper.GetString("kafka.topic.nonce"),
			},
			broker.NewEngineConsumerGroupHandler(engine.GlobalEngine()),
		)
	})
}
