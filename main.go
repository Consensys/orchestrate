package main

import (
	"context"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"

	commonHandlers "gitlab.com/ConsenSys/client/fr/core-stack/common.git/handlers"
	core "gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	types "gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
	hand "gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-decoder.git/handlers"
)

// TxDecoder is the handler used by the Sarama consumer of the tx-decoder worker
type TxDecoder struct {
	w              *core.Worker
	saramaProducer sarama.SyncProducer
	mec            *ethclient.MultiEthClient
	cfg            Config
}

// Setup configure handler
func (h *TxDecoder) Setup(s sarama.ConsumerGroupSession) error {
	// Instantiate workers
	h.w = core.NewWorker(h.cfg.Worker.Slots)

	// Worker::logger
	h.w.Use(hand.LoggerHandler)

	// Worker::unmarchaller
	h.w.Use(commonHandlers.Loader(infSarama.NewUnmarshaller()))

	// Worker::marker
	h.w.Use(commonHandlers.Marker(infSarama.NewSimpleOffsetMarker(s)))

	// Worker::decoder
	registry := hand.LoadABIRegistry(h.cfg.App.ABIs)
	h.w.Use(
		hand.TransactionDecoder(registry),
	)

	// Worker::producer
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
		commonHandlers.Producer(
			infSarama.NewProducer(
				h.saramaProducer,
				prepareMsg,
			),
		),
	)

	return nil
}

// ConsumeClaim consume messages from queue
func (h *TxDecoder) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
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
func (h *TxDecoder) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

func main() {
	// Load Config from env variables
	var cfg Config
	LoadConfig(&cfg)

	// Configure the logger
	ConfigureLogger(cfg.Log)
	log.Info("Initialize worker...")

	// Init config
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	config.Consumer.Return.Errors = true
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	// Create sarama client
	client, err := sarama.NewClient(cfg.Kafka.Address, config)
	if err != nil {
		log.Fatalf("Could not to start sarama client: %v", err)
		return
	}
	defer client.Close()
	var brokers = make(map[int32]string)
	for _, v := range client.Brokers() {
		brokers[v.ID()] = v.Addr()
	}
	log.WithFields(log.Fields{
		"kafka.endpoint": cfg.Kafka.Address,
		"kafka.brokers":  brokers,
	}).Info("Kafka client ready")

	// Create sarama sync producer
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Fatalf("Could not start sarama producer: %v", err)
		return
	}
	log.Info("Kafka producer ready")
	defer p.Close()

	// Create sarama consumer
	g, err := sarama.NewConsumerGroupFromClient(cfg.Kafka.ConsumerGroup, client)
	if err != nil {
		log.Fatalf("Could not start sarama consumer group: %v", err)
		return
	}
	coordinatorBroker, _ := client.Coordinator(cfg.Kafka.ConsumerGroup)
	log.WithFields(log.Fields{
		"kafka.consumergroup": cfg.Kafka.ConsumerGroup,
		"kafka.coordinator":   coordinatorBroker.Addr(),
	}).Info("Kafka consumer group ready")
	defer func() { g.Close() }()

	// Create an ethereum client connection
	mec, err := ethclient.MultiDial(cfg.Eth.URLs)
	if err != nil {
		log.Fatalf("Could not to start ethereum client: %v", err)
	}

	// Listen to multi in-topics depending on the chainID listened by tx-listener
	var multiChainInTopics []string
	for _, chainID := range mec.Networks(context.Background()) {
		log.WithFields(log.Fields{
			"ethclient.chainID": chainID,
		}).Info("Ethereum client ready")
		multiChainInTopics = append(multiChainInTopics, cfg.Kafka.InTopic+"-"+chainID.String())
	}

	txDecoder := &TxDecoder{mec: mec, saramaProducer: p, cfg: cfg}
	log.WithFields(log.Fields{
		"kafka.topics": multiChainInTopics,
	}).Info("Starting worker")
	g.Consume(context.Background(), multiChainInTopics, txDecoder)
}
