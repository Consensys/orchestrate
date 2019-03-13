package main

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	listener "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/core/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

// Logger to log context elements before and after the worker
func Logger(ctx *core.Context) {
	log.WithFields(log.Fields{
		"Chain":       ctx.T.Chain.Id,
		"BlockNumber": ctx.T.Receipt.BlockNumber,
		"TxIndex":     ctx.T.Receipt.TxIndex,
		"TxHash":      ctx.T.Receipt.TxHash,
	}).Debug("tx-listener-worker: new receipt")

	ctx.Next()

	if len(ctx.T.Errors) != 0 {
		log.WithFields(log.Fields{
			"Chain":  ctx.T.Chain.Id,
			"TxHash": ctx.T.Receipt.TxHash,
		}).Errorf("tx-listener-worker: Errors: %v", ctx.T.Errors)
	}
}

// ReceiptUnmarshaller assumes that input message is a go-ethereum receipt
type ReceiptUnmarshaller struct{}

// Unmarshal message expected to be a trace protobuffer
func (u *ReceiptUnmarshaller) Unmarshal(msg interface{}, t *trace.Trace) error {
	// Cast message into receipt
	receipt, ok := msg.(*listener.TxListenerReceipt)
	if !ok {
		return fmt.Errorf("Message does not match expected format")
	}

	t.Chain = (&common.Chain{}).SetID(receipt.ChainID)
	t.Receipt = ethereum.FromGethReceipt(&receipt.Receipt).
		SetBlockHash(receipt.BlockHash).
		SetBlockNumber(uint64(receipt.BlockNumber)).
		SetTxHash(receipt.TxHash).
		SetTxIndex(receipt.TxIndex)

	return nil
}

// Loader creates an handler loading input
func Loader(u services.Unmarshaller) core.HandlerFunc {
	return func(ctx *core.Context) {
		// Unmarshal message
		err := u.Unmarshal(ctx.Msg, ctx.T)
		if err != nil {
			// TODO: handle error
			ctx.AbortWithError(err)
			return
		}
	}
}

// TxListenerHandler is an handler consuming receipts
type TxListenerHandler struct {
	w *core.Worker
}

// Setup configure the handler
func (h *TxListenerHandler) Setup() error {
	// Instantiate worker
	h.w = core.NewWorker(1)

	// Handler::loader
	h.w.Use(Loader(&ReceiptUnmarshaller{}))

	// Handler::logger
	h.w.Use(Logger)

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
		log.Info("Worker channel closed...")
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

func main() {
	// Create an ethereum client connection
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	ethURLs := []string{
		"https://ropsten.infura.io/v3/81e039ce6c8a465180822b525e3644d7",
		"https://rinkeby.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://kovan.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
		"https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
	}
	log.Infof("Connecting to EthClients: %v", ethURLs)
	mec, err := ethclient.MultiDial(ethURLs)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("Multi-Eth client ready")

	// Initialize listener configuration
	listenerCfg := listener.NewConfig()
	listenerCfg.TxListener.Return.Blocks = true
	listenerCfg.TxListener.Return.Errors = true

	viper.Set("blockcursor.backoff", time.Second)
	viper.Set("blockcursor.limit", 40)

	txlistener := listener.NewTxListener(listener.NewEthClient(mec), listenerCfg)

	// Create and Listener Handler
	handler := TxListenerHandler{}
	handler.Setup()
	log.Infof("Worker ready")

	// Start listening all chains
	for _, chainID := range mec.Networks(context.Background()) {
		txlistener.Listen(chainID, -1, 0)
	}

	// Start listening
	handler.Listen(txlistener)
}
