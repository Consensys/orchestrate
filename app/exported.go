package app

import (
	"context"
	"fmt"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/loader"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/offset"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	server "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/tests/e2e.git/handlers"
)

var (
	app       *App
	startOnce = &sync.Once{}
)

func init() {
	// Create app
	app = NewApp()

	// Set Kafka Group value
	viper.Set("kafka.group", "group-e2e")
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
	wg := sync.WaitGroup{}

	// Initialize Engine
	wg.Add(1)
	go func() {
		engine.Init(ctx)
		wg.Done()
	}()

	// Initialize Handlers
	wg.Add(1)
	go func() {
		handlers.Init(ctx)
		wg.Done()
	}()

	// Initialize ConsumerGroup
	wg.Add(1)
	go func() {
		// Want to consume from the oldest offset
		broker.InitConfig()
		config := broker.GlobalConfig()
		config.Consumer.Offsets.Initial = sarama.OffsetOldest

		broker.InitConsumerGroup(ctx)
		wg.Done()
	}()

	// Initialize Ethereum client
	wg.Add(1)
	go func() {
		rpc.Init(ctx)
		wg.Done()
	}()

	// Wait for engine and handlers to be ready
	wg.Wait()
}

func registerHandlers() {
	wg := sync.WaitGroup{}

	// Register handlers on engine
	wg.Add(1)
	go func() {
		// Generic handlers on every worker
		engine.Register(logger.Logger)
		engine.Register(loader.Loader)
		engine.Register(offset.Marker)
		wg.Done()
	}()

	// Wait for ConsumerGroup & Engine to be ready
	wg.Wait()
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
		initComponents(cancelCtx)

		// Register all Handlers
		registerHandlers()

		// Indicate that application is ready
		app.ready.Store(true)

		// Send Transactions
		SendTx()

		// // Start consuming on every topics
		// // Initialize Topics list by chain
		topics := []string{
			viper.GetString("kafka.topic.crafter"),
			viper.GetString("kafka.topic.nonce"),
			viper.GetString("kafka.topic.signer"),
			viper.GetString("kafka.topic.sender"),
			viper.GetString("kafka.topic.decoded"),
		}
		for _, chainID := range rpc.GlobalClient().Networks(context.Background()) {
			topics = append(topics, fmt.Sprintf("%v-%v", viper.GetString("kafka.topic.decoder"), chainID.String()))
		}

		log.Debugf("worker: start consuming on %q", topics)
		err := broker.Consume(
			cancelCtx,
			topics,
			broker.NewEngineConsumerGroupHandler(engine.GlobalEngine()),
		)
		if err != nil {
			log.WithError(err).Error("worker: error on consumer")
		}
	})
}
