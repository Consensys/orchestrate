package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	encoding "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/producer"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/ethereum/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/ethereum/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/utils"
	handlercfg "gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/handler/base"
	handler "gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/handler/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/listener"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/listener/base"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
)

// Loader is a Middleware enginer.HandlerFunc that Load sarama.ConsumerGroup messages
func Loader(txctx *engine.TxContext) {
	// Cast message into sarama.ConsumerMessage
	receipt, ok := txctx.In.(*types.TxListenerReceipt)
	if !ok {
		txctx.Logger.Errorf("loader: expected a types.TxListenerReceipt")
		_ = txctx.AbortWithError(fmt.Errorf("invalid input message format"))
		return
	}

	// Set receipt
	txctx.Envelope.Receipt = ethereum.FromGethReceipt(&receipt.Receipt).
		SetBlockHash(receipt.BlockHash).
		SetBlockNumber(uint64(receipt.BlockNumber)).
		SetTxIndex(receipt.TxIndex)
	txctx.Envelope.Chain = (&chain.Chain{}).SetID(receipt.ChainID)

	// Enrich Logger
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"chain.id":     receipt.ChainID.Text(10),
		"tx.hash":      receipt.TxHash.Hex(),
		"tx.index":     receipt.TxIndex,
		"block.hash":   receipt.BlockHash.Hex(),
		"block.Number": uint64(receipt.BlockNumber),
	})

	txctx.Logger.Tracef("loader: message loaded: %v", txctx.Envelope.String())
}

// PrepareMsg prepare message to produce from TxContexts
func PrepareMsg(txctx *engine.TxContext, msg *sarama.ProducerMessage) error {
	// // Marshal Envelope into sarama Message
	err := encoding.Marshal(txctx.Envelope, msg)
	if err != nil {
		return err
	}

	// Set Topic to Nonce topic
	msg.Topic = utils.KafkaChainTopic(viper.GetString("kafka.topic.decoder"), txctx.Envelope.GetChain().ID())

	log.WithFields(log.Fields{
		"topic": msg.Topic,
	}).Infof("message prepared")

	return nil
}

// Producer creates a producer handler
func Producer(p sarama.SyncProducer) engine.HandlerFunc {
	return producer.Producer(p, PrepareMsg)
}

func main() {
	// Set log configuration
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

	// Initialize Listener
	viper.Set("eth.clients", []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		// "https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		// "https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		// "https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	})
	viper.Set("kafka.topic.decoder", "dodo")

	// Initialize listener config than initialize listener
	config := base.NewConfig()
	config.TxListener.Return.Blocks = true
	config.TxListener.Return.Errors = true
	listener.SetGlobalConfig(config)
	listener.Init(context.Background())

	// Initialize engine and register handlers
	engine.Init(context.Background())

	// Init Producer
	broker.InitSyncProducer(context.Background())

	engine.Register(logger.Logger)
	engine.Register(Loader)
	engine.Register(Producer(broker.GlobalSyncProducer()))

	// Create handler
	conf, err := handlercfg.NewConfig()
	if err != nil {
		log.WithError(err).Fatalf("listener: could not create config")
	}
	h := handler.NewHandler(engine.GlobalEngine(), broker.GlobalClient(), broker.GlobalSyncProducer(), conf)

	ctx, cancel := context.WithCancel(context.Background())
	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { cancel() })
	defer sig.Close()

	// Start listening
	err = listener.Listen(
		ctx,
		[]*big.Int{
			// big.NewInt(1),
			big.NewInt(3),
		},
		h,
	)
	if err != nil {
		log.WithError(err).Error("listener: error listening")
	}
}
