package txcrafter

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/viper"
	chaininjector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/chain-injector"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/crafter"
	gasestimator "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/gas/gas-estimator"
	gaspricer "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/gas/gas-pricer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/loader/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"
	nonceattributor "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/nonce/attributor"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/offset"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/opentracing"
	producer "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/producer/tx-crafter"
	injector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/trace-injector"
	txupdater "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/tx_updater"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
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
		// Initialize trace injector
		func() {
			ctxWithValue := context.WithValue(ctx, serviceName("service-name"), viper.GetString(jaeger.ServiceNameViperKey))
			injector.Init(ctxWithValue)
		},
		// Initialize Multi-tenancy
		func() {
			multitenancy.Init(ctx)
		},

		// Initialize crafter
		func() {
			crafter.Init(ctx)
		},

		// Initialize Gas Estimator
		func() {
			gasestimator.Init(ctx)
		},

		// Initialize Gas Pricer
		func() {
			gaspricer.Init(ctx)
		},

		// Initialize Nonce Attributor
		func() {
			nonceattributor.Init(ctx)
		},

		// Initialize Updater
		func() {
			txupdater.Init()
		},

		// Initialize Producer
		func() {
			producer.Init(ctx)
		},

		// Initialize GetBigChainID injector
		func() {
			chaininjector.Init(ctx)
		},
	)
}

func initComponents(ctx context.Context) {
	utils.InParallel(
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
			// Set Kafka Group value
			viper.Set(broker.KafkaGroupViperKey, "group-crafter")
			broker.InitConsumerGroup(ctx)
		},
	)
}

func registerHandlers() {
	// Register handlers on engine
	// Generic handlers on every worker
	engine.Register(opentracing.GlobalHandler())
	engine.Register(logger.Logger("info"))
	engine.Register(sarama.Loader)
	engine.Register(offset.Marker)
	engine.Register(opentracing.GlobalHandler())
	engine.Register(producer.GlobalHandler())
	engine.Register(txupdater.GlobalHandler())
	engine.Register(injector.GlobalHandler())
	engine.Register(multitenancy.GlobalHandler())

	// Specific handlers tk Tx-Crafter worker
	engine.Register(chaininjector.GlobalHandler())
	engine.Register(crafter.GlobalHandler())
	engine.Register(gaspricer.GlobalHandler())
	engine.Register(gasestimator.GlobalHandler())
	engine.Register(nonceattributor.GlobalHandler())
}

// Start starts application
func Run(ctx context.Context) error {
	var err error
	runOnce.Do(func() {
		// Register all Handlers
		initComponents(ctx)
		registerHandlers()

		topics := []string{
			viper.GetString(broker.TxCrafterViperKey),
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
