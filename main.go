package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	handCom "gitlab.com/ConsenSys/client/fr/core-stack/common.git/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/trace"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener"
	infSarama "gitlab.com/ConsenSys/client/fr/core-stack/infra/sarama.git"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/handlers"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/infra"
)

// TxListenerHandler is an handler consuming receipts
type TxListenerHandler struct {
	w   *core.Worker
	p   sarama.SyncProducer
	cfg Config
}

// Setup configure the handler
func (h *TxListenerHandler) Setup() error {
	// Instantiate worker
	h.w = core.NewWorker(1)

	// Handler::loader
	h.w.Use(handCom.Loader(&infra.ReceiptUnmarshaller{}))

	// Handler::logger
	h.w.Use(handlers.Logger)

	// Handler::Producer
	marshaller := infSarama.NewMarshaller()

	prepareMsg := func(t *types.Trace, msg *sarama.ProducerMessage) error {
		err := marshaller.Marshal(t, msg)
		if err != nil {
			return err
		}

		// Set topic
		msg.Topic = fmt.Sprintf("%v-%v", h.cfg.Kafka.OutTopic, t.Chain().ID.Text(16))
		return nil
	}

	h.w.Use(
		handCom.Producer(
			infSarama.NewProducer(
				h.p,
				prepareMsg,
			),
		),
	)

	return nil
}

// Listen start listening to receipts from listener
func (h *TxListenerHandler) Listen(l listener.TxListener) error {
	go func() {
		for err := range l.Errors() {
			log.WithFields(log.Fields{
				"Chain": err.ChainID.Text(16),
			}).Errorf("tx-listener-worker: error %v", err)
		}
	}()

	go func() {
		for block := range l.Blocks() {
			log.WithFields(log.Fields{
				"BlockHash":   block.Header().Hash().Hex(),
				"BlockNumber": block.Header().Number,
				"Chain":       block.ChainID.Text(16),
			}).Debugf("tx-listener-worker: new block")
		}
	}()

	in := make(chan interface{})
	go func() {
		// Pipe channels for interface compatibility
		for msg := range l.Receipts() {
			in <- msg
		}
		close(in)
	}()
	log.Info("Start worker...")
	h.w.Run(in)

	return nil
}

// Cleanup cleans handler
func (h *TxListenerHandler) Cleanup() error {
	return nil
}

func newTxListener(cfg Config, mec *ethclient.MultiEthClient) (listener.TxListener, error) {

	return listener.NewTxListener(mec), nil
}

func getLastMessage(client sarama.Client, topic string, partition int32) (*sarama.Record, error) {
	// Retrieve last offset that has been produced
	lastOffset, err := client.GetOffset(topic, partition, -1)
	if err != nil {
		return nil, err
	}

	// Retrieve last offset that has been produced
	broker, err := client.Leader(topic, partition)
	if err != nil {
		return nil, err
	}

	req := &sarama.FetchRequest{
		MinBytes:    client.Config().Consumer.Fetch.Min,
		MaxWaitTime: int32(client.Config().Consumer.MaxWaitTime / time.Millisecond),
	}
	req.AddBlock(topic, 0, lastOffset-1, client.Config().Consumer.Fetch.Max)
	req.Version = 4
	req.Isolation = sarama.ReadUncommitted
	response, err := broker.Fetch(req)

	if err != nil {
		return nil, err
	}
	block := response.GetBlock(topic, partition)
	if len(block.RecordsSet) == 0 {
		return nil, nil
	}
	records := block.RecordsSet[0]
	record := records.RecordBatch.Records[0]
	return record, nil
}

func main() {
	// Load Config from env variables
	var cfg Config
	LoadConfig(&cfg)

	// Configure the logger
	ConfigureLogger(cfg.Log)

	// Init Sarama config
	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V1_0_0_0
	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Producer.Return.Errors = true
	saramaCfg.Producer.Return.Successes = true

	// Create sarama client
	client, err := sarama.NewClient(cfg.Kafka.Address, saramaCfg)
	if err != nil {
		log.Error(err)
		return
	}
	defer client.Close()
	log.Infof("Sarama client ready")

	// Create sarama sync producer
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Sarama producer ready")
	defer p.Close()

	// Create an ethereum client connection
	log.Infof("Connecting to EthClients: %v", cfg.Eth.URLs)
	mec, err := ethclient.MultiDial(cfg.Eth.URLs)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Multi-Eth client ready")

	// Create and Listener Handler
	handler := TxListenerHandler{p: p, cfg: cfg}
	handler.Setup()
	log.Infof("Worker ready")

	// Initialize listener configuration
	listenerCfg := listener.NewConfig()
	backoff, err := time.ParseDuration(cfg.Listener.Block.Backoff)
	if err != nil {
		log.Error(err)
		return
	}
	listenerCfg.BlockCursor.Backoff = backoff
	listenerCfg.BlockCursor.Limit = cfg.Listener.Block.Limit
	listenerCfg.TxListener.Return.Blocks = true
	listenerCfg.TxListener.Return.Errors = true

	txlistener := listener.NewTxListener(listener.NewEthClient(mec, listenerCfg))

	// Start listening all chains
	for _, chainID := range mec.Networks(context.Background()) {
		var blockNumber, txIndex int64
		// Determine starting position
		if pos, ok := cfg.Listener.Start.Specific[hexutil.EncodeBig(chainID)]; ok {
			blockNumber, txIndex, err = ParseStartingPosition(pos)
			if err != nil {
				log.Error(err)
				return
			}
		} else {
			blockNumber, err = TranslateBlockNumber(cfg.Listener.Start.Default)
			if err != nil {
				log.Error(err)
				return
			}
		}

		if blockNumber == -2 {
			// Retrieve last Offset for every chain
			outTopic := chainTopic(cfg.Kafka.OutTopic, chainID)

			lastRecord, err := getLastMessage(client, outTopic, 0)
			if err != nil {
				log.Error(err)
				return
			}

			if lastRecord == nil {
				blockNumber, txIndex = 0, 0
			} else {
				var pb tracepb.Trace
				// Unmarshal record to protobuffer
				err = proto.Unmarshal(lastRecord.Value, &pb)
				if err != nil {
					log.Error(err)
					return
				}
				blockNumber, txIndex = int64(pb.Receipt.BlockNumber), int64(pb.Receipt.TxIndex)+1
			}
		}

		// Start listening
		txlistener.Listen(chainID, blockNumber, txIndex, listenerCfg)
	}

	// Start listening
	handler.Listen(txlistener)
	"os"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/boilerplate-worker.git/cmd"
)

func main() {
	command := cmd.NewCommand()

	if err := command.Execute(); err != nil {
		log.Errorf("%v\n", err)
		os.Exit(1)
	}
}
