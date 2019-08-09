package app

import (
	"context"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/common"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/loader"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/offset"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/opentracing"
	server "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-sender.git/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-sender.git/handlers/producer"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-sender.git/handlers/sender"
)

var (
	app       *App
	startOnce = &sync.Once{}
)

var l = createLogger("worker")

func init() {
	// Create app
	app = NewApp()

	// Set Kafka Group value
	viper.Set("kafka.group", "group-sender")
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
		// Initialize Engine
		func() { engine.Init(ctx) },

		// Initialize Handlers
		func() { handlers.Init(ctx) },

		// Initialize ConsumerGroup
		func() { broker.InitConsumerGroup(ctx) },
	)
}

func registerHandlers() {
	// Generic handlers on every worker
	engine.Register(logger.Logger)
	engine.Register(loader.Loader)
	engine.Register(offset.Marker)
	engine.Register(producer.GlobalHandler())
	engine.Register(opentracing.GlobalHandler())

	// Specific handlers tk Sender worker
	engine.Register(sender.GlobalHandler())
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

		topics := []string{
			viper.GetString("kafka.topic.sender"),
		}
		l.WithFields(log.Fields{
			"topics": topics,
		}).Info("connecting")
		// Start consuming on topic tx-sender
		err := broker.Consume(
			cancelCtx,
			topics,
			broker.NewEngineConsumerGroupHandler(engine.GlobalEngine()),
		)
		if err != nil {
			l.WithError(err).Error("error on consumer")
		}
	})
}

func createLogger(name string) *log.Entry {
	return log.WithField("name", name)
}
