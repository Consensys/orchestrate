package txdecoder

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/decoder"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/loader/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/offset"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/opentracing"
	producer "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/producer/tx-decoder"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	server "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/http/healthcheck"
)

var (
	app       *common.App
	startOnce = &sync.Once{}
)

func init() {
	// Create app
	app = common.NewApp()
}

func startServer(ctx context.Context) {
	// Initialize server
	server.Init(ctx)

	// Register Healthcheck
	server.Enhance(healthcheck.HealthCheck(app))

	// Start Listening
	_ = server.ListenAndServe()
}

type serviceName string

func initHandlers(ctx context.Context) {
	common.InParallel(
		// Initialize Jaeger tracer
		func() {
			ctx = context.WithValue(ctx, serviceName("service-name"), viper.GetString("jaeger.service.name"))
			opentracing.Init(ctx)
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

func initConsumerGroup(ctx context.Context) {
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
			viper.Set("kafka.group", "group-decoder")
			broker.InitConsumerGroup(ctx)
		},
		// Initialize Ethereum client
		func() {
			rpc.Init(ctx)
		},
	)

	// Generic handlers on every worker
	engine.Register(logger.Logger)
	engine.Register(sarama.Loader)
	engine.Register(offset.Marker)
	engine.Register(producer.GlobalHandler())
	engine.Register(opentracing.GlobalHandler())

	// Specific handlers of Tx-Decoder worker
	engine.Register(decoder.GlobalHandler())
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

		// Initialize ConsumerGroup
		initConsumerGroup(cancelCtx)

		// Initialize Topics list by chain
		var topics []string
		for _, chainID := range rpc.GlobalClient().Networks(context.Background()) {
			topics = append(topics, fmt.Sprintf("%v-%v", viper.GetString("kafka.topic.decoder"), chainID.String()))
		}

		log.Infof("Connecting to the following Kafka topics %v", topics)

		// Indicate that application is ready
		// TODO: we need to update so ready can append when Consume has finished to Setup
		app.SetReady(true)

		// Start consuming on topic tx-decoder
		err := broker.Consume(
			cancelCtx,
			topics,
			broker.NewEngineConsumerGroupHandler(engine.GlobalEngine()),
		)
		if err != nil {
			log.WithError(err).Error("worker: failed to consume messages")
		}
	})
}
