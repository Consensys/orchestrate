package main

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	listener "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/encoding"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

// Logger to log context elements before and after the worker
func Logger(txctx *engine.TxContext) {
	log.WithFields(log.Fields{
		"Chain":       txctx.Envelope.Chain.Id,
		"BlockNumber": txctx.Envelope.Receipt.BlockNumber,
		"TxIndex":     txctx.Envelope.Receipt.TxIndex,
		"TxHash":      txctx.Envelope.Receipt.TxHash,
	}).Debug("tx-listener: new receipt")

	txctx.Next()

	if len(txctx.Envelope.Errors) != 0 {
		log.WithFields(log.Fields{
			"Chain":  txctx.Envelope.Chain.Id,
			"TxHash": txctx.Envelope.Receipt.TxHash,
		}).Errorf("tx-listener: Errors: %v", txctx.Envelope.Errors)
	}
}

// ReceiptUnmarshaller assumes that input message is a go-ethereum receipt
type ReceiptUnmarshaller struct{}

// Unmarshal message expected to be a Envelope protobuffer
func (u *ReceiptUnmarshaller) Unmarshal(msg interface{}, e *envelope.Envelope) error {
	// Cast message into receipt
	receipt, ok := msg.(*listener.TxListenerReceipt)
	if !ok {
		return fmt.Errorf("message does not match expected format")
	}

	e.Chain = (&common.Chain{}).SetID(receipt.ChainID)
	e.Receipt = ethereum.FromGethReceipt(&receipt.Receipt).
		SetBlockHash(receipt.BlockHash).
		SetBlockNumber(uint64(receipt.BlockNumber)).
		SetTxHash(receipt.TxHash).
		SetTxIndex(receipt.TxIndex)

	return nil
}

// Loader creates an handler loading input
func Loader(u encoding.Unmarshaller) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// Unmarshal message
		err := u.Unmarshal(txctx.Msg, txctx.Envelope)
		if err != nil {
			// TODO: handle error
			_ = txctx.AbortWithError(err)
			return
		}
	}
}

// TxListenerHandler is an handler consuming receipts
type TxListenerHandler struct {
	engine *engine.Engine
}

// Setup configure the handler
func (h *TxListenerHandler) Setup() {
	// Instantiate engine
	cfg := engine.NewConfig()
	h.engine = engine.NewEngine(&cfg)

	// Handler::loader
	h.engine.Register(Loader(&ReceiptUnmarshaller{}))

	// Handler::logger
	h.engine.Register(Logger)
}

// Listen start listening to receipts from listener
func (h *TxListenerHandler) Listen(l listener.TxListener) {
	go func() {
		for err := range l.Errors() {
			log.WithFields(log.Fields{
				"Chain": err.ChainID.Text(16),
			}).Errorf("tx-listener: error %v", err)
		}
	}()

	go func() {
		for block := range l.Blocks() {
			log.WithFields(log.Fields{
				"BlockHash":   block.Header().Hash().Hex(),
				"BlockNumber": block.Header().Number,
				"Chain":       block.ChainID.Text(16),
			}).Debugf("tx-listener: new block")
		}
	}()

	in := make(chan interface{})
	go func() {
		// Pipe channels for interface compatibility
		for msg := range l.Receipts() {
			in <- msg
		}
		log.Info("Engine channel closed...")
		close(in)
	}()
	log.Info("Start engine...")
	h.engine.Run(context.Background(), in)
}

// Cleanup cleans handler
func (h *TxListenerHandler) Cleanup() error {
	return nil
}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	viper.Set("eth.clients", []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		"https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	})

	// Initialize Multi-Client
	ethclient.Init(context.Background())

	// Initialize listener configuration
	config := listener.NewConfig()
	config.TxListener.Return.Blocks = true
	config.TxListener.Return.Errors = true

	viper.Set("blockcursor.backoff", time.Second)
	viper.Set("blockcursor.limit", 40)

	txlistener := listener.NewTxListener(ethclient.GlobalMultiClient(), config)

	// Create and Listener Handler
	handler := TxListenerHandler{}
	handler.Setup()
	log.Infof("Engine ready")

	// Start listening all chains
	for _, chainID := range ethclient.GlobalMultiClient().Networks(context.Background()) {
		_, _ = txlistener.Listen(chainID, -1, 0)
	}

	// Start listening
	handler.Listen(txlistener)
}
