package main

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"

	commonHandlers "gitlab.com/ConsenSys/client/fr/core-stack/common.git/handlers"
	core "gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	types "gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	infEth "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
	hand "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers"
)

var opts Config

type handler struct {
	w *core.Worker
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

func prepareMsg(t *types.Trace, msg *sarama.ProducerMessage) error {
	marshaller := infSarama.NewMarshaller()

	err := marshaller.Marshal(t, msg)
	if err != nil {
		return err
	}

	// Set topic
	msg.Topic = opts.Kafka.OutTopic
	return nil
}

func newEthClient(rawurl string) *infEth.EthClient {
	ec, err := infEth.Dial(rawurl)
	if err != nil {
		panic(err)
	}
	log.Info("Connected to Ethereum client")
	return ec
}

// Setup configure handler
func (h *handler) Setup(s sarama.ConsumerGroupSession) error {
	// Instantiate workers
	h.w = core.NewWorker(opts.App.WorkerSlots)

	// Worker::logger
	h.w.Use(hand.LoggerHandler)

	// Worker::unmarchaller
	h.w.Use(commonHandlers.Loader(infSarama.NewUnmarshaller()))

	// Worker::marker
	h.w.Use(commonHandlers.Marker(infSarama.NewSimpleOffsetMarker(s)))

	// Worker::signer
	txSigner := infEth.NewStaticSigner(opts.App.Vault.Accounts)
	h.w.Use(
		hand.Signer(txSigner),
	)

	// Worker::producer
	h.w.Use(
		commonHandlers.Producer(
			infSarama.NewProducer(
				newSaramaSyncProducer(newSaramaClient([]string{opts.Kafka.Address})),
				prepareMsg,
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
	log.Info("Start worker...")

	client := newSaramaClient([]string{opts.Kafka.Address})

	// Create consumer
	g, err := sarama.NewConsumerGroupFromClient(opts.Kafka.ConsumerGroup, client)
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

	g.Consume(context.Background(), []string{opts.Kafka.InTopic}, &handler{})
}
