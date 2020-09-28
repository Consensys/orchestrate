package txsender

import (
	"context"
	"sync"
	"time"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/viper"
	chaininjector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/chain-injector"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/loader/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"
	noncechecker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/nonce/checker"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/offset"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/opentracing"
	producer "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/producer/tx-sender"
	rawdecoder "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/raw-decoder"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/sender"
	injector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/trace-injector"
	txupdater "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/tx_updater"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

var (
	appli   *app.App
	runOnce = &sync.Once{}
)

type serviceName string

func initHandlers(ctx context.Context) {
	utils.InParallel(
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
		// Initialize injector manager
		func() { injector.Init(ctx) },
		// Initialize Tx Updater
		func() { txupdater.Init() },
		// Initialize producer
		func() { producer.Init(ctx) },
		// Initialize GetBigChainID injector
		func() {
			chaininjector.Init(ctx)
		},
		func() {
			viper.Set(utils.RetryMaxIntervalViperKey, 1*time.Second)
			viper.Set(utils.RetryMaxElapsedTimeViperKey, 15*time.Second)
			ethclient.Init(ctx)
		},
	)
}

func initComponents(ctx context.Context) {
	utils.InParallel(
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
	engine.Register(txupdater.GlobalHandler())
	engine.Register(chaininjector.GlobalHandler())
	engine.Register(rawdecoder.RawDecoder)
	engine.Register(noncechecker.GlobalRecoveryStatusSetter())
	engine.Register(injector.GlobalHandler())
	engine.Register(noncechecker.GlobalChecker())
	engine.Register(sender.GlobalHandler())
}

// Run starts application
func Run(ctx context.Context) error {
	var err error
	runOnce.Do(func() {
		// Register all Handlers
		initComponents(ctx)
		registerHandlers()

		topics := []string{
			viper.GetString(broker.TxSenderViperKey),
		}

		// Create appli to expose metrics
		appli, err = New(
			app.NewConfig(viper.GetViper()),
			broker.NewConsumerDaemon(
				broker.GlobalClient(),
				broker.GlobalSyncProducer(),
				broker.GlobalConsumerGroup(),
				topics,
				broker.NewEngineConsumerGroupHandler(engine.GlobalEngine()),
			),
		)
		if err != nil {
			log.FromContext(ctx).WithError(err).Info("could not create app")
			return
		}

		err = appli.Run(ctx)
	})

	return err
}
