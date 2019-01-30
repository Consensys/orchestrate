package main

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	handCom "gitlab.com/ConsenSys/client/fr/core-stack/common.git/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	infEth "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
	hand "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-signer.git/handlers"
)

// TxSignerHandler is the handler used by the Sarama consumer of the tx-signer worker
type TxSignerHandler struct {
	w              *core.Worker
	saramaProducer sarama.SyncProducer
	mec            *ethclient.MultiEthClient
	cfg            Config
}

// Setup configure handler
func (h *TxSignerHandler) Setup(s sarama.ConsumerGroupSession) error {
	// Instantiate workers
	h.w = core.NewWorker(h.cfg.Worker.Slots)

	// Worker::unmarchaller
	h.w.Use(handCom.Loader(infSarama.NewUnmarshaller()))

	// Worker::logger
	h.w.Use(hand.Logger)

	// Worker::marker
	h.w.Use(handCom.Marker(infSarama.NewSimpleOffsetMarker(s)))

	// Worker::signer
	txSigner := infEth.NewStaticSigner(h.cfg.Vault.Accounts)
	h.w.Use(
		hand.Signer(txSigner),
	)

	// Handler::Producer
	marshaller := infSarama.NewMarshaller()

	prepareMsg := func(t *types.Trace, msg *sarama.ProducerMessage) error {
		err := marshaller.Marshal(t, msg)
		if err != nil {
			return err
		}

		// Set topic
		msg.Topic = h.cfg.Kafka.OutTopic
		return nil
	}

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
func (h *TxSignerHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
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
func (h *TxSignerHandler) Cleanup(s sarama.ConsumerGroupSession) error {
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
	// Load Config from env variables
	var cfg Config
	LoadConfig(&cfg)

	// Configure the logger
	ConfigureLogger(cfg.Log)
	log.Info("Start worker...")

	// Init config
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	// Create sarama client
	client, err := sarama.NewClient(cfg.Kafka.Address, config)
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Close()
	log.Println("Client ready")

	// Create sarama sync producer
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Producer ready")
	defer p.Close()

	// Create sarama consumer
	g, err := sarama.NewConsumerGroupFromClient(cfg.Kafka.ConsumerGroup, client)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Consumer group ready")
	defer func() { g.Close() }()

	// Create an ethereum client connection
	mec, err := ethclient.MutiDial(cfg.Eth.URLs)
	if err != nil {
		log.Errorf("Got error %v", err)
	}

	txCrafter := &TxSignerHandler{mec: mec, saramaProducer: p, cfg: cfg}
	err = g.Consume(context.Background(), []string{cfg.Kafka.InTopic}, txCrafter)
	log.Error(err)
}
