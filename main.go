package main

import (
	"context"
	"math/big"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	handCom "gitlab.com/ConsenSys/client/fr/core-stack/common.git/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	infEth "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
	hand "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/handlers"
	infra "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/infra"
)

// TxCrafter is the handler used by the Sarama consumer of the tx-craft worker
type TxCrafter struct {
	w              *core.Worker
	saramaProducer sarama.SyncProducer
	mec            *ethclient.MultiEthClient
	cfg            Config
}

// Setup configure the handler
func (h *TxCrafter) Setup(s sarama.ConsumerGroupSession) error {
	// Instantiate worker
	h.w = core.NewWorker(h.cfg.Worker.Slots)

	// Handler::loader
	h.w.Use(handCom.Loader(infSarama.NewUnmarshaller()))

	// Handler::logger
	h.w.Use(hand.Logger)

	// Handler::marker
	h.w.Use(handCom.Marker(infSarama.NewSimpleOffsetMarker(s)))

	// Handler::Faucet
	crediter, err := infra.NewSaramaCrediter(h.cfg.Faucet, h.saramaProducer)
	if err != nil {
		return err
	}
	faucet, err := infra.CreateFaucet(h.cfg.Faucet, h.mec.PendingBalanceAt, crediter.Credit)
	if err != nil {
		return err
	}
	creditAmount := big.NewInt(0)
	creditAmount.SetString(h.cfg.Faucet.CreditAmount, 10)
	h.w.Use(hand.Faucet(faucet, creditAmount))

	// Handler::Crafter
	crafter := infEth.PayloadCrafter{}
	registry := infra.NewERC1400ABIRegistry()
	h.w.Use(hand.Crafter(registry, &crafter))

	// Handler:Gas
	gasManager := infEth.NewGasManager(h.mec)
	h.w.Use(hand.GasPricer(gasManager))    // Gas Price
	h.w.Use(hand.GasEstimator(gasManager)) // Gas Limit

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
func (h *TxCrafter) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
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
func (h *TxCrafter) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
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

	txCrafter := &TxCrafter{mec: mec, saramaProducer: p, cfg: cfg}
	err = g.Consume(context.Background(), []string{cfg.Kafka.InTopic}, txCrafter)
	log.Error(err)
}
