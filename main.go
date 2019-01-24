package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	handCom "gitlab.com/ConsenSys/client/fr/core-stack/common.git/handlers"
	core "gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	infEth "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
	hand "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/handlers"
	inf "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-crafter.git/infra"
)

// TxCrafter is the handler used by the Sarama consumer of the tx-craft worker
type TxCrafter struct {
	w              *core.Worker
	saramaProducer sarama.SyncProducer
	ethClient      *infEth.EthClient
	cfg            Config
}

func (h *TxCrafter) prepareMsg(t *types.Trace, msg *sarama.ProducerMessage) error {
	marshaller := infSarama.NewMarshaller()

	err := marshaller.Marshal(t, msg)
	if err != nil {
		return err
	}

	// Set topic
	msg.Topic = h.cfg.Kafka.TopicOut
	return nil
}

// Setup configure the handler
func (h *TxCrafter) Setup(s sarama.ConsumerGroupSession) error {
	// Instantiate workers
	h.w = core.NewWorker(50)

	// TODO : to be removed ?
	// Worker::logger middleware which will log before and after the tx-crafter nonce action
	h.w.Use(hand.Logger)

	// Sarama message unmarshalling loader
	h.w.Use(handCom.Loader(infSarama.NewUnmarshaller()))

	// Crafter
	crafter := infEth.PayloadCrafter{}
	registry := hand.NewERC1400ABIRegistry()
	h.w.Use(hand.Crafter(registry, &crafter))

	gasManager := infEth.NewSimpleGasManager(h.ethClient)
	h.w.Use(hand.GasPricer(gasManager))    // Gas Price
	h.w.Use(hand.GasEstimator(gasManager)) // Gas Limit

	// Faucet
	faucetAddress := common.HexToAddress(h.cfg.Faucet.Address)
	faucetTopicOut := h.cfg.Kafka.TopicIn
	faucet, err := inf.CreateFaucet(
		h.cfg.Eth.URL,
		faucetAddress,
		h.cfg.Faucet.BalanceMax,
		h.saramaProducer,
		faucetTopicOut,
	)
	if err != nil {
		return err
	}
	faucetAmount := big.NewInt(h.cfg.Faucet.Amount)
	h.w.Use(hand.Faucet(faucet, faucetAmount))

	// Sarama producer
	msgProducer := infSarama.NewProducer(h.saramaProducer, h.prepareMsg)
	h.w.Use(handCom.Producer(msgProducer))

	// Sarama marker
	h.w.Use(handCom.Marker(infSarama.NewSimpleOffsetMarker(s)))

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
	client, err := sarama.NewClient([]string{cfg.Kafka.Address}, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()
	fmt.Println("Client ready")

	// Create sarama sync producer
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Producer ready")
	defer p.Close()

	// Create sarama consumer
	g, err := sarama.NewConsumerGroupFromClient(cfg.Kafka.ConsumerGroup, client)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("Consumer Group ready")
	defer func() { g.Close() }()

	// Create an ethereum client connection
	e, err := infEth.Dial(cfg.Eth.URL)
	if err != nil {
		log.Errorf("Got error %v", err)
	}

	txCrafter := &TxCrafter{ethClient: e, saramaProducer: p}
	g.Consume(context.Background(), []string{cfg.Kafka.TopicIn}, txCrafter)
}
