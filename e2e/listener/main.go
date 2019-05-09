package main

import (
	"context"
	"fmt"
	"math/big"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	handlercfg "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener/handler/base"
	handler "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener/handler/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener/listener"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener/listener/base"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/types"
	broker "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/handlers/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

// Loader is a Middleware enginer.HandlerFunc that Load sarama.ConsumerGroup messages
func Loader(txctx *engine.TxContext) {
	// Cast message into sarama.ConsumerMessage
	receipt, ok := txctx.Msg.(*types.TxListenerReceipt)
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
	txctx.Envelope.Chain = (&common.Chain{}).SetID(receipt.ChainID)

	// Enrich Logger
	txctx.Logger = txctx.Logger.WithFields(log.Fields{
		"chain.id":   receipt.ChainID.Text(10),
		"tx.hash":    receipt.TxHash.Hex(),
		"block.hash": receipt.BlockHash.Hex(),
	})

	txctx.Logger.Tracef("loader: message loaded: %v", txctx.Envelope.String())
}

func main() {
	// Set log configuration
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

	// Initialize Listener
	viper.Set("eth.clients", []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		"https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	})

	// Initialize listener config than initialize listener
	config := base.NewConfig()
	config.TxListener.Return.Blocks = true
	config.TxListener.Return.Errors = true
	listener.SetGlobalConfig(config)
	listener.Init(context.Background())

	// Initialize engine and register handlers
	engine.Init(context.Background())
	engine.Register(Loader)
	engine.Register(logger.Logger)

	// Init Producer
	broker.InitSyncProducer(context.Background())

	// Create handler
	conf, err := handlercfg.NewConfig()
	if err != nil {
		log.WithError(err).Fatalf("listener: could not create config")
	}
	h := handler.NewHandler(engine.GlobalEngine(), broker.GlobalClient(), broker.GlobalSyncProducer(), conf)

	// Start listening
	_ = listener.Listen(
		context.Background(),
		[]*big.Int{
			big.NewInt(1),
			big.NewInt(3),
		},
		h,
	)
}
