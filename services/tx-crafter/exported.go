package txcrafter

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	chaininjector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/chain-injector"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/crafter"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/faucet"
	gasestimator "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/gas/gas-estimator"
	gaspricer "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/gas/gas-pricer"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/loader/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"
	nonceattributor "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/nonce/attributor"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/offset"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/opentracing"
	producer "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/producer/tx-crafter"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/tessera"
	injector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/trace-injector"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing/jaeger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

var (
	appli     *app.App
	startOnce = &sync.Once{}
	done      chan struct{}
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

		// Initialize faucet
		func() {
			faucet.Init(ctx)
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
		// Initialize Producer
		func() {
			producer.Init(ctx)
		},

		// Initialize GetBigChainID injector
		func() {
			chaininjector.Init(ctx)
		},

		// Initialize Tessera client
		func() {
			tessera.Init(ctx)
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
	engine.Register(injector.GlobalHandler())
	engine.Register(multitenancy.GlobalHandler())

	// Specific handlers tk Tx-Crafter worker
	engine.Register(chaininjector.GlobalHandler())
	engine.Register(faucet.GlobalHandler())
	engine.Register(crafter.GlobalHandler())
	engine.Register(gaspricer.GlobalHandler())
	engine.Register(gasestimator.GlobalHandler())
	engine.Register(nonceattributor.GlobalHandler())
	engine.Register(tessera.GlobalHandler())
}

// Start starts application
func Start(ctx context.Context) (chan struct{}, error) {
	var err error
	startOnce.Do(func() {
		// Chan to notify that sub-go routines stopped
		done = make(chan struct{})

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
			viper.GetString(broker.TxCrafterViperKey),
		}

		go func() {
			log.FromContext(ctx).WithFields(logrus.Fields{
				"topics": topics,
			}).Info("connecting")

			errConsume := broker.Consume(
				ctx,
				topics,
				broker.NewEngineConsumerGroupHandler(engine.GlobalEngine()),
			)
			if errConsume != nil {
				log.FromContext(ctx).WithError(errConsume).Error("error on consumer")
			}
			close(done)
		}()
	})
	return done, err
}

func Stop(ctx context.Context) error {
	<-done
	return errors.CombineErrors(broker.Stop(ctx), appli.Stop(ctx))
}
