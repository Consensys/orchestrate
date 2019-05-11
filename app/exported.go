package app

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient/rpc"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/loader"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/offset"
	server "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-decoder.git/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-decoder.git/handlers/decoder"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-decoder.git/handlers/producer"
)

var (
	app       *App
	startOnce = &sync.Once{}
)

func init() {
	// Create app
	app = NewApp()

	// Set Kafka Group value
	viper.Set("kafka.group", "group-decoder")
}

func startServer(ctx context.Context) {
	// Initialize server
	server.Init(ctx)

	// Register Healthcheck
	server.Enhance(healthcheck.HealthCheck(app))

	// Start Listening
	_ = server.ListenAndServe()
}

func initConsumerGroup(ctx context.Context) {
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

	// Register handlers on engine
	wg.Add(1)
	go func() {
		// Generic handlers on every worker
		engine.Register(logger.Logger)
		engine.Register(loader.Loader)
		engine.Register(offset.Marker)
		engine.Register(producer.GlobalHandler())

		// Specific handlers of Tx-Decoder worker
		engine.Register(decoder.GlobalHandler())
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
		initConsumerGroup(cancelCtx)

		// Initialize Topics list by chain
		topics := []string{}
		for _, chainID := range rpc.GlobalClient().Networks(context.Background()) {
			topics = append(topics, fmt.Sprintf("%v-%v", viper.GetString("kafka.topic.decoder"), chainID.String()))
		}

		// Indicate that application is ready
		// TODO: we need to update so ready can append when Consume has finished to Setup
		app.ready.Store(true)

		// Start consuming on topic tx-decoder
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
