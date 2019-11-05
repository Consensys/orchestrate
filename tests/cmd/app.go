package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	loader "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/loader/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/handlers/offset"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	server "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/handlers/dispatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/tests/service/cucumber/steps"
)

var (
	app         *common.App
	readyToTest chan bool
	startOnce   = &sync.Once{}
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

func LongKeyOf(topics map[string]string) dispatcher.KeyOfFunc {
	return func(txtcx *engine.TxContext) (string, error) {
		topic, ok := topics[txtcx.In.Entrypoint()]
		if !ok {
			return "", fmt.Errorf("unknown message entrypoint")
		}

		scenario, ok := txtcx.Envelope.GetMetadataValue("scenario.id")
		if !ok {
			return "", fmt.Errorf("message has no test scenario")
		}

		return steps.LongKeyOf(topic, scenario, txtcx.Envelope.GetMetadata().Id), nil
	}
}

func ShortKeyOf(topics map[string]string) dispatcher.KeyOfFunc {
	return func(txtcx *engine.TxContext) (string, error) {
		topic, ok := topics[txtcx.In.Entrypoint()]
		if !ok {
			return "", fmt.Errorf("unknown message entrypoint")
		}

		scenario, ok := txtcx.Envelope.GetMetadataValue("scenario.id")
		if !ok {
			return "", fmt.Errorf("message has no test scenario")
		}

		return steps.ShortKeyOf(topic, scenario), nil
	}
}

func initComponents(ctx context.Context) {
	common.InParallel(
		// Initialize Engine
		func() {
			engine.Init(ctx)
		},
		// Initialize ConsumerGroup
		func() {
			viper.Set(broker.KafkaGroupViperKey, "group-e2e")
			broker.InitConsumerGroup(ctx)
		},
		// Initialize Handlers
		func() {
			// Prepare topics map for dispatcher
			topics := make(map[string]string)
			for _, topic := range steps.TOPICS {
				topics[viper.GetString(fmt.Sprintf("topic.%v", topic))] = topic
			}
			dispatcher.SetKeyOfFuncs(
				LongKeyOf(topics),
				ShortKeyOf(topics),
			)
			handlers.Init(ctx)
		},
		// Initialize cucumber handlers
		func() {
			cucumber.Init(ctx)
		},
	)
}

func registerHandlers() {
	// Generic handlers on every worker
	engine.Register(loader.Loader)
	engine.Register(offset.Marker)
	engine.Register(dispatcher.GlobalHandler())
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
		app.SetReady(true)

		// Start consuming on every topics of interest
		var topics []string
		for _, topic := range steps.TOPICS {
			topics = append(topics, viper.GetString(fmt.Sprintf("topic.%v", topic)))
		}

		readyToTest = make(chan bool, 1)

		go func() {
			<-readyToTest
			cucumber.Run(cancel, cucumber.GlobalOptions())
		}()

		cg := &EmbeddingConsumerGroupHandler{
			engine: broker.NewEngineConsumerGroupHandler(engine.GlobalEngine()),
		}

		log.Debugf("worker: start consuming on %q", topics)
		err := broker.Consume(
			cancelCtx,
			topics,
			cg,
		)
		if err != nil {
			log.WithError(err).Fatalf("worker: error on consumer with topics: %s", topics)
		}

	})
}

type EmbeddingConsumerGroupHandler struct {
	engine *broker.EngineConsumerGroupHandler
}

func (h *EmbeddingConsumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error {
	err := h.engine.Setup(s)
	readyToTest <- true
	return err
}

func (h *EmbeddingConsumerGroupHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	return h.engine.ConsumeClaim(s, c)
}

func (h *EmbeddingConsumerGroupHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	return h.engine.Cleanup(s)
}
