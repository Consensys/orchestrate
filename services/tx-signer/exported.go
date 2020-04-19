package txsigner

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/loader/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/offset"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/opentracing"
	producer "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/producer/tx-signer"
	injector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/trace-injector"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/vault"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app/worker"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	orchlog "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

var (
	appli     *app.App
	startOnce = &sync.Once{}
	done      chan struct{}
	cancel    func()
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
			viper.Set(broker.KafkaGroupViperKey, "group-signer")
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
	engine.Register(opentracing.GlobalHandler())
	engine.Register(producer.GlobalHandler())
	engine.Register(injector.GlobalHandler())
	engine.Register(multitenancy.GlobalHandler())

	// Specific handlers for Signer worker
	engine.Register(vault.GlobalHandler())
}

// Start starts application
func Start(ctx context.Context) error {
	var err error
	startOnce.Do(func() {
		// Create Configuration
		cfg := app.NewConfig(viper.GetViper())
		orchlog.ConfigureLogger(cfg.HTTP)

		ctx, cancel = context.WithCancel(ctx)

		// Register all Handlers
		initComponents(ctx)
		registerHandlers()

		// Create appli to expose metrics
		appli, err = worker.New(cfg, prom.DefaultRegisterer)
		if err != nil {
			return
		}

		err = appli.Start(ctx)
		if err != nil {
			return
		}

		// Start consuming on topic tx-sender
		topics := []string{
			viper.GetString(broker.TxSignerViperKey),
			viper.GetString(broker.AccountGeneratorViperKey),
		}

		done = make(chan struct{})
		go func() {
			log.FromContext(ctx).WithFields(logrus.Fields{
				"topics": topics,
			}).Info("connecting")

			err = broker.Consume(
				ctx,
				topics,
				broker.NewEngineConsumerGroupHandler(engine.GlobalEngine()),
			)
			if err != nil {
				log.FromContext(ctx).WithError(err).Error("error on consumer")
			}
			close(done)
		}()
	})

	return nil
}

func Stop(ctx context.Context) error {
	cancel()
	err := appli.Stop(ctx)
	<-done
	return err
}
