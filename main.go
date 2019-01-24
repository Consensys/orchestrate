package main

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"

	handCom "gitlab.com/ConsenSys/client/fr/core-stack/common.git/handlers"
	core "gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	infEth "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
	infRedis "gitlab.com/ConsenSys/client/fr/core-stack/infra/redis.git"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
	hand "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-nonce.git/handlers"
)

var opts Config

// TxNonceHandler is the handler used by the Sarama consumer of the tx-nonce worker
type TxNonceHandler struct {
	w              *core.Worker
	saramaProducer sarama.SyncProducer
	ethClient      *infEth.EthClient
}

func prepareMsg(t *types.Trace, msg *sarama.ProducerMessage) error {
	marshaller := infSarama.NewMarshaller()

	err := marshaller.Marshal(t, msg)
	if err != nil {
		return err
	}

	// Set topic
	msg.Topic = opts.App.OutTopic
	return nil
}

// Setup configure the handler
func (h *TxNonceHandler) Setup(s sarama.ConsumerGroupSession) error {
	// Instantiate workers
	h.w = core.NewWorker(opts.App.WorkerSlots)

	// Worker::logger
	h.w.Use(hand.LoggerHandler)

	// Worker::marker
	h.w.Use(handCom.Marker(infSarama.NewSimpleOffsetMarker(s)))

	// Worker::nonce
	h.w.Use(
		hand.NonceHandler(
			infRedis.NewNonceManager(opts.Conn.Redis.URL, opts.Conn.Redis.LockTimeout),
			hand.GetChainNonce(h.ethClient),
		),
	)

	// Worker::producer
	h.w.Use(
		handCom.Producer(
			infSarama.NewProducer(
				h.saramaProducer,
				prepareMsg,
			),
		),
	)

	return nil
}

// ConsumeClaim consume messages from queue
func (h *TxNonceHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	in := make(chan interface{})
	go func() {
		// Pipe channels for interface compatibility
		for msg := range c.Messages() {
			in <- msg
		}
		close(in)
	}()
	h.w.Run(in)

	return nil
}

// Cleanup cleans handler
func (h *TxNonceHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

func main() {
	LoadConfig(&opts)
	ConfigureLogger(opts.Log)

	// Init config
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	// Create sarama client
	client, err := sarama.NewClient([]string{opts.Conn.Kafka.URL}, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()
	log.Info("Sarama client ready")

	// Create sarama sync producer
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Info("Producer ready")
	defer p.Close()

	// Create sarama consumer
	g, err := sarama.NewConsumerGroupFromClient(opts.App.ConsumerGroup, client)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Consumer Group ready")
	defer func() { g.Close() }()

	// Create an ethereum client connection
	e, err := infEth.Dial(opts.Conn.ETHClient.URL)
	if err != nil {
		// TODO log with logger from worker
		log.Errorf("Got error %v", err)
	}

	// Track errors
	go func() {
		for err := range g.Errors() {
			log.Error("ERROR", err)
		}
	}()

	txNonceHandler := &TxNonceHandler{ethClient: e, saramaProducer: p}
	g.Consume(context.Background(), []string{opts.App.InTopic}, txNonceHandler)
}
