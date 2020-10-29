package txsigner

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/loader/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/offset"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/opentracing"
	producer "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/producer/tx-signer"
	injector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/trace-injector"
	txupdater "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/tx_updater"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault"
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
		// Initialize Jaeger tracer injector
		func() {
			ctxWithValue := context.WithValue(ctx, serviceName("service-name"), viper.GetString(jaeger.ServiceNameViperKey))
			injector.Init(ctxWithValue)
		},
		// Initialize Multi-tenancy
		func() {
			multitenancy.Init(ctx)
		},
		// Initialize Vault
		func() { vault.Init(ctx) },
		// Initialize Tx Updater
		func() { txupdater.Init() },
		// Initialize Producer
		func() { producer.Init(ctx) },
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
			broker.InitConsumerGroup(ctx, "group-signer")
		},
	)
}

func registerHandlers() {
	// Generic handlers on every worker
	engine.Register(opentracing.GlobalHandler())
	engine.Register(logger.Logger("info"))
	engine.Register(sarama.Loader)
	engine.Register(offset.Marker)
	engine.Register(opentracing.GlobalHandler())
	engine.Register(producer.GlobalHandler())
	engine.Register(txupdater.GlobalHandler())
	engine.Register(injector.GlobalHandler())

	// Specific handlers for Signer worker
	engine.Register(vault.GlobalHandler())
}

// Run starts application
func Run(ctx context.Context) error {
	var err error
	runOnce.Do(func() {
		// Register all Handlers
		initComponents(ctx)
		registerHandlers()

		topics := []string{
			viper.GetString(broker.TxSignerViperKey),
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
