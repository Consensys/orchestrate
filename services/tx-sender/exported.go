package txsender

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
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
	engine.Register(chaininjector.GlobalHandler())
	engine.Register(rawdecoder.RawDecoder)
	engine.Register(noncechecker.GlobalRecoveryStatusSetter())
	engine.Register(injector.GlobalHandler())
	engine.Register(noncechecker.GlobalChecker())

	engine.Register(sender.GlobalHandler())
}

// Start starts application
func Start(ctx context.Context) error {
	var err error
	startOnce.Do(func() {
		ctx, cancel = context.WithCancel(ctx)

		// Register all Handlers
		initComponents(ctx)
		registerHandlers()

		// Create appli to expose metrics
		appli, err = app.New(
			app.NewConfig(viper.GetViper()),
			app.MetricsOpt(),
		)
		if err != nil {
			return
		}

		err = appli.Start(ctx)
		if err != nil {
			return
		}

		// Start consuming on topic tx-sender
		topics := []string{
			viper.GetString(broker.TxSenderViperKey),
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

	return err
}

func Stop(ctx context.Context) error {
	cancel()
	err := appli.Stop(ctx)
	<-done
	return err
}
