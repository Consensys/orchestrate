package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/boilerplate-worker.git/cmd"
	handCom "gitlab.com/ConsenSys/client/fr/core-stack/common.git/handlers"
	core "gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	infEth "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
	hand "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-sender.git/handlers"
)

// TxSenderHandler is SaramaHandler used in tx-nonce worker
type TxSenderHandler struct {
	w   *core.Worker
	mec *ethclient.MultiEthClient
	cfg Config
}

// Setup configure the handler
func (h *TxSenderHandler) Setup(s sarama.ConsumerGroupSession) error {
	// Instantiate worker
	h.w = core.NewWorker(h.cfg.Worker.Slots)

	// Handler::loader
	h.w.Use(handCom.Loader(infSarama.NewUnmarshaller()))

	// Handler::logger
	h.w.Use(hand.Logger)

	// Handler::marker
	h.w.Use(handCom.Marker(infSarama.NewSimpleOffsetMarker(s)))

	// Handler::Sender
	h.w.Use(hand.Sender(infEth.NewTxSender(h.mec)))

	return nil
}

// ConsumeClaim consume messages from queue
func (h *TxSenderHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
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
func (h *TxSenderHandler) Cleanup(s sarama.ConsumerGroupSession) error {
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

	// Create sarama client
	client, err := sarama.NewClient(cfg.Kafka.Address, sc)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()
	log.Info("Sarama client ready")

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
	log.Info("Multi-Client ready")

	// Track errors
	go func() {
		for err := range g.Errors() {
			log.Error("ERROR", err)
		}
	}()

	g.Consume(
		context.Background(),
		[]string{cfg.Kafka.InTopic},
		&TxSenderHandler{mec: mec, cfg: cfg},
	)
	command := cmd.NewCommand()

	if err := command.Execute(); err != nil {
		log.Errorf("%v\n", err)
		os.Exit(1)
	}
}
