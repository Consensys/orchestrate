package main

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	handCom "gitlab.com/ConsenSys/client/fr/core-stack/common.git/handlers"
	core "gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	infRedis "gitlab.com/ConsenSys/client/fr/core-stack/infra/redis.git"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
	hand "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-nonce.git/handlers"
)

// TxNonceHandler is SaramaHandler used in tx-nonce worker
type TxNonceHandler struct {
	w              *core.Worker
	saramaProducer sarama.SyncProducer
	mec            *ethclient.MultiEthClient
	cfg            Config
}

func (h *TxNonceHandler) prepareMsg(t *types.Trace, msg *sarama.ProducerMessage) error {
	marshaller := infSarama.NewMarshaller()

	err := marshaller.Marshal(t, msg)
	if err != nil {
		return err
	}

	// Set topic
	msg.Topic = h.cfg.Kafka.OutTopic
	return nil
}

// Setup configure the handler
func (h *TxNonceHandler) Setup(s sarama.ConsumerGroupSession) error {
	// Instantiate worker
	h.w = core.NewWorker(h.cfg.Worker.Slots)

	// Hanlder::loader
	h.w.Use(handCom.Loader(infSarama.NewUnmarshaller()))

	// Hanlder::logger
	h.w.Use(hand.Logger)

	// Hanlder::marker
	h.w.Use(handCom.Marker(infSarama.NewSimpleOffsetMarker(s)))

	// Hanlder::nonce
	h.w.Use(
		hand.NonceHandler(
			infRedis.NewNonceManager(h.cfg.Redis.Address, h.cfg.Redis.LockTimeout),
			h.mec.PendingNonceAt,
		),
	)

	// Hanlder::producer
	h.w.Use(
		handCom.Producer(
			infSarama.NewProducer(
				h.saramaProducer,
				h.prepareMsg,
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
	var cfg Config
	LoadConfig(&cfg)
	ConfigureLogger(cfg.Log)

	// Init Sarama config
	sc := sarama.NewConfig()
	sc.Version = sarama.V1_0_0_0
	sc.Consumer.Return.Errors = true
	sc.Producer.Return.Errors = true
	sc.Producer.Return.Successes = true

	// Create sarama client
	client, err := sarama.NewClient(cfg.Kafka.Address, sc)
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
	g, err := sarama.NewConsumerGroupFromClient(cfg.Kafka.ConsumerGroup, client)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Consumer Group ready")
	defer func() { g.Close() }()

	// Create a Multi Ethereum Client connection
	mec, err := ethclient.MutiDial(cfg.Eth.URLs)
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

	g.Consume(
		context.Background(),
		[]string{cfg.Kafka.InTopic},
		&TxNonceHandler{mec: mec, saramaProducer: p, cfg: cfg},
	)
}
