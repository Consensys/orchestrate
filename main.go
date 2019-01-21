package main

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"

	coreHandlers "gitlab.com/ConsenSys/client/fr/core-stack/core.git/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/infra"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/trace"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	nonceHandlers "gitlab.com/ConsenSys/client/fr/core-stack/worker/nonce/handlers"
)

var opts Config

type handler struct {
	w *types.Worker
}

func newSaramaSyncProducer(client sarama.Client) sarama.SyncProducer {
	// Create producer
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	log.Info("Producer ready")
	return p
}

func ctxToProducerMessage(pb *tracepb.Trace) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic:     opts.App.OutTopic,
		Partition: -1,
	}
	b, _ := proto.Marshal(pb)
	msg.Value = sarama.ByteEncoder(b)
	return msg
}

func newEthClient(rawurl string) *infra.EthClient {
	ec, err := infra.Dial(rawurl)
	if err != nil {
		panic(err)
	}
	log.Info("Connected to Ethereum client")
	return ec
}

// Setup configure handler
func (h *handler) Setup(s sarama.ConsumerGroupSession) error {
	// Instantiate workers
	h.w = types.NewWorker(50)

	// Worker::logger
	h.w.Use(nonceHandlers.LoggerHandler)

	// Worker::marker
	h.w.Use(coreHandlers.Marker(infra.NewSimpleSaramaOffsetMarker(s)))

	// Worker::nonce
	h.w.Use(
		coreHandlers.NonceHandler(
			infra.NewRedisNonceManager(opts.Conn.Redis.URL),
			coreHandlers.GetChainNonce(newEthClient(opts.Conn.ETHClient.URL)),
		),
	)

	// Worker::producer
	h.w.Use(
		coreHandlers.Producer(
			infra.NewSaramaProducer(
				newSaramaSyncProducer(newSaramaClient([]string{opts.Conn.Kafka.URL})),
				ctxToProducerMessage,
			),
		),
	)

	return nil
}

// ConsumeClaim consume messages from queue
func (h *handler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
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
func (h *handler) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

func newSaramaClient(kafkaURL []string) sarama.Client {
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	// Create client
	client, err := sarama.NewClient(kafkaURL, config)
	if err != nil {
		panic(err)
	}
	log.Info("Sarama client ready")
	return client
}

func main() {
	LoadConfig(&opts)
	ConfigureLogger(opts.Log)

	client := newSaramaClient([]string{opts.Conn.Kafka.URL})

	// Create consumer
	g, err := sarama.NewConsumerGroupFromClient(opts.App.ConsumerGroup, client)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Consumer Group ready")
	defer func() { g.Close() }()

	// Track errors
	go func() {
		for err := range g.Errors() {
			log.Error("ERROR", err)
		}
	}()

	g.Consume(context.Background(), []string{opts.App.InTopic}, &handler{})
}
