package base

import (
	"context"
	"math/big"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/handler"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/engine"
)

// Handler implements TxListenerHandler interface
//
// It uses a pkg Engine to listen to chains messages
type Handler struct {
	engine *engine.Engine

	Conf *Config
}

// NewHandler creates a new EngineConsumerGroupHandler
func NewHandler(e *engine.Engine, conf *Config) *Handler {
	return &Handler{
		engine: e,
		Conf:   conf,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (h *Handler) Setup(s handler.TxListenerSession) error {
	chains := []string{}
	for _, chain := range s.Chains() {
		chains = append(chains, chain.Text(10))
	}

	log.WithFields(log.Fields{
		"chains": chains,
	}).Info("tx-listener: ready to listen to chains")

	return nil
}

// Listen starts Listening.
// Once the Messages() channel is closed it finishes its processing and exits loop
//
// Make sure that you have registered the chain of HandlerFunc on context before ConsumeClaim is called
func (h *Handler) Listen(session handler.TxListenerSession, l handler.ChainListener) error {
	logger := log.WithFields(log.Fields{
		"chain.id": l.ChainID().Text(10),
	})

	block, txIdx := l.InitialPosition()

	// Start listener in a separate goroutine
	logger.WithFields(log.Fields{
		"start.block":    block,
		"start.tx-index": txIdx,
	}).Infof("tx-listener: start listening")

	// Attach ConsumerGroupSession to context
	wait := &sync.WaitGroup{}
	wait.Add(3)

	// Start consuming blocks
	go func() {
	blockLoop:
		for {
			select {
			case block, ok := <-l.Blocks():
				if !ok {
					break blockLoop
				}
				logger.WithFields(log.Fields{
					"block.hash":   block.Header().Hash().Hex(),
					"block.number": block.Header().Number,
				}).Debugf("tx-listener: new block")
			case <-l.Context().Done():
				break blockLoop
			}
		}
		wait.Done()
	}()

	// Start consuming errors
	go func() {
	errorLoop:
		for {
			select {
			case err, ok := <-l.Errors():
				if !ok {
					break errorLoop
				}
				logger.WithError(err).Debugf("tx-listener: error")
			case <-l.Context().Done():
				break errorLoop
			}
		}
		wait.Done()
	}()

	// Start consuming Receipts
	go func() {
		h.engine.Run(l.Context(), Pipe(l.Context(), l.Receipts()))
		wait.Done()
	}()

	wait.Wait()
	logger.Infof("tx-listener: stoped listening")

	return nil
}

// GetInitialPosition return initial position
func (h *Handler) GetInitialPosition(chain *big.Int) (blockNumber, txIndex int64, err error) {
	position, ok := h.Conf.Start.Positions[chain.Text(10)]
	if !ok {
		blockNumber = h.Conf.Start.Default.BlockNumber
		txIndex = h.Conf.Start.Default.TxIndex
	} else {
		blockNumber = position.BlockNumber
		txIndex = position.TxIndex
	}

	if blockNumber >= -1 {
		return blockNumber, txIndex, nil
	}

	return -1, 0, nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
// but before the offsets are committed for the very last time.
func (h *Handler) Cleanup(session handler.TxListenerSession) error {
	h.engine.CleanUp()
	log.Infof("tx-listener: cleaned up")
	return nil
}

// Pipe take a channel of types.TxListenerReceipt and pipes it into a channel of interface{}
//
// Pipe will stop forwarding messages either
// - receipt channel is closed
// - ctx has been canceled
func Pipe(ctx context.Context, receiptChan <-chan *types.TxListenerReceipt) <-chan interface{} {
	interfaceChan := make(chan interface{})

	// Start a goroutine that pipe messages
	go func() {
	pipeLoop:
		for {
			select {
			case msg, ok := <-receiptChan:
				if !ok {
					// Sarama channel has been closed so we exit loop
					break pipeLoop
				}
				interfaceChan <- msg
			case <-ctx.Done():
				// Context has been cancel so we exit loop
				break pipeLoop
			}
		}
		close(interfaceChan)
	}()

	return interfaceChan
}
