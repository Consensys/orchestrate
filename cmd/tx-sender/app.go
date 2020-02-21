package txsender

import (
	"context"
	"sync"

	rawdecoder "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/raw-decoder"

	chaininjector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/chain-injector"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/loader/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/logger"
	noncechecker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/nonce/checker"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/offset"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/opentracing"
	producer "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/producer/tx-sender"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/sender"
	injector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/trace-injector"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing/jaeger"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/key"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
)

var (
	app       = common.NewApp()
	startOnce = &sync.Once{}
)

var l = createLogger("worker")

type serviceName string

func initHandlers(ctx context.Context) {
	common.InParallel(
		// Initialize Jaeger tracer
		func() {
			ctxWithValue := context.WithValue(ctx, serviceName("service-name"), viper.GetString(jaeger.ServiceNameViperKey))
			opentracing.Init(ctxWithValue)
		},
		// Initialize Jaeger tracer injector
		func() {
			ctxWithValue := context.WithValue(ctx, serviceName("service-name"), viper.GetString(jaeger.ServiceNameViperKey))
			injector.Init(ctxWithValue)
		},
		// Initialize Multi-tenancy
		func() {
			multitenancy.Init(ctx)
		},
		// Initialize sender
		func() { sender.Init(ctx) },
		// Initialize nonce manager
		func() { noncechecker.Init(ctx) },
		// Initialize producer
		func() { producer.Init(ctx) },
		// Initialize GetBigChainID injector
		func() {
			chaininjector.Init(ctx)
		},
	)
}

func initComponents(ctx context.Context) {
	common.InParallel(
		// Initialize Engine
		func() { engine.Init(ctx) },

		// Initialize Handlers
		func() { initHandlers(ctx) },

		// Initialize ConsumerGroup
		func() {
			// Set Kafka Group value
			viper.Set(broker.KafkaGroupViperKey, "group-sender")
			broker.InitConsumerGroup(ctx)
		},
	)
}

func registerHandlers() {
	// Generic handlers on every worker
	engine.Register(opentracing.GlobalHandler())
	engine.Register(logger.Logger("info"))
	engine.Register(sarama.Loader)
	engine.Register(offset.Marker)
	engine.Register(multitenancy.GlobalHandler())
	engine.Register(opentracing.GlobalHandler())

	// Recovery Status Setter surrounds the producer
	// c.f. docstring RecoveryStatusSetter handler
	engine.Register(producer.GlobalHandler())
	engine.Register(chaininjector.GlobalHandler())
	engine.Register(rawdecoder.RawDecoder)
	engine.Register(noncechecker.GlobalRecoveryStatusSetter())
	engine.Register(injector.GlobalHandler())
	engine.Register(noncechecker.GlobalChecker())

	engine.Register(sender.GlobalHandler())
}

// Start starts application
func Start(ctx context.Context) {
	startOnce.Do(func() {
		apiKey := viper.GetString(authkey.APIKeyViperKey)
		if apiKey != "" {
			// Inject authorization header in context for later authentication
			ctx = authutils.WithAPIKey(ctx, apiKey)
		}

		cancelCtx, cancel := context.WithCancel(ctx)
		go metrics.StartServer(ctx, cancel, app.IsAlive, app.IsReady)

		// Initialize all components of the server
		initComponents(cancelCtx)

		// Register all Handlers
		registerHandlers()

		// Indicate that application is ready
		// TODO: we need to update so SetReady can be called when Consume has finished to Setup
		app.SetReady(true)

		topics := []string{
			viper.GetString(broker.TxSenderViperKey),
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
