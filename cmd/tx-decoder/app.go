package txdecoder

import (
	"context"
	"fmt"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/multitenancy"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/metrics"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/decoder"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/loader/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/offset"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/opentracing"
	producer "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/producer/tx-decoder"
	injector "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/trace-injector"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing/jaeger"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/key"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
)

var (
	app       = common.NewApp()
	startOnce = &sync.Once{}
)

type serviceName string

func initHandlers(ctx context.Context) {
	common.InParallel(
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
		// Initialize decoder
		func() {
			decoder.Init(ctx)
		},
		// Initialize Producer
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
			// Set Kafka Group value
			viper.Set(broker.KafkaGroupViperKey, "group-decoder")
			broker.InitConsumerGroup(ctx)
		},
		// Initialize Ethereum client
		func() {
			rpc.Init(ctx)
		},
	)

	// Generic handlers on every worker
	engine.Register(opentracing.GlobalHandler())
	engine.Register(logger.Logger("debug"))
	engine.Register(sarama.Loader)
	engine.Register(offset.Marker)
	engine.Register(opentracing.GlobalHandler())
	engine.Register(producer.GlobalHandler())
	engine.Register(injector.GlobalHandler())
	engine.Register(multitenancy.GlobalHandler())

	// Specific handlers of Tx-Decoder worker
	engine.Register(decoder.GlobalHandler())
}

// Start starts application
func Start(ctx context.Context) {
	startOnce.Do(func() {
		apiKey := viper.GetString(authkey.APIKeyViperKey)
		if apiKey != "" {
			// Inject authorization header in context for later authentication
			ctx = authutils.WithAuthorization(ctx, fmt.Sprintf("APIKey %v", apiKey))
		}

		cancelCtx, cancel := context.WithCancel(ctx)
		go metrics.StartServer(ctx, cancel, app.IsAlive, app.IsReady)

		// Initialize ConsumerGroup
		initComponents(cancelCtx)

		// Indicate that application is ready
		// TODO: we need to update so SetReady can be called when Consume has finished to Setup
		app.SetReady(true)

		// Start consuming on topic tx-decoder
		err := broker.Consume(
			cancelCtx,
			[]string{viper.GetString(broker.TxDecoderViperKey)},
			broker.NewEngineConsumerGroupHandler(engine.GlobalEngine()),
		)
		if err != nil {
			log.WithError(err).Error("worker: failed to consume messages")
		}
	})
}
