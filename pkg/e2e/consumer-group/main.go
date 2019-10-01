package main

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/handlers/offset"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/utils"
)

// Define a handler method
func handler(txctx *engine.TxContext) {
	txctx.Logger.WithFields(log.Fields{
		"offset":    txctx.In.(*broker.Msg).Offset,
		"partition": txctx.In.(*broker.Msg).Partition,
	}).Infof("Handling message")

	// Simulate latency
	time.Sleep(10 * time.Millisecond)
}

func main() {
	// Initialize consumer group
	viper.Set("worker.group", "e2e-group")
	broker.InitConsumerGroup(context.Background())

	// Initialize engine
	engine.Init(context.Background())

	// Register Marker Middleware handler
	engine.Register(offset.Marker)

	// Register Pipeline handler
	engine.Register(handler)

	// Create ConsumerGroupHandler
	groupHandler := broker.NewEngineConsumerGroupHandler(engine.GlobalEngine())

	// Listen to Signals for graceful close purpose
	ctx, cancel := context.WithCancel(context.Background())
	utils.NewSignalListener(func(s os.Signal) { cancel() })

	// Start consuming
	_ = broker.Consume(
		ctx,
		[]string{"topic-e2e"},
		groupHandler,
	)

	log.Infof("Gracefully closed")
}
