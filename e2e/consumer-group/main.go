package main

import (
	"context"
	"os"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/offset"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
)

// Define a handler method
func handler(txctx *engine.TxContext) {
	txctx.Logger.WithFields(log.Fields{
		"offset":    txctx.Msg.(*sarama.ConsumerMessage).Offset,
		"partition": txctx.Msg.(*sarama.ConsumerMessage).Partition,
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
