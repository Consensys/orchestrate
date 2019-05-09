package base

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	cursor "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener/block-cursor/base"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener/handler"
	tiptracker "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener/tip-tracker/base"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/types"
)

// Client interface for a TxListener
type Client interface {
	ethclient.ChainLedgerReader
	ethclient.ChainSyncReader
}

type TxListener struct {
	conf *Config

	ec Client

	wait      *sync.WaitGroup
	closeOnce *sync.Once
	closed    chan struct{}
}

// NewTxListener creates a new Listener
func NewTxListener(ec Client, conf *Config) *TxListener {
	return &TxListener{
		conf:      conf,
		ec:        ec,
		wait:      &sync.WaitGroup{},
		closeOnce: &sync.Once{},
		closed:    make(chan struct{}),
	}
}

func (l *TxListener) Listen(ctx context.Context, chains []*big.Int, h handler.TxListenerHandler) error {
	select {
	case <-l.closed:
		return fmt.Errorf("listener closed")
	default:
	}

	// Start new listener session
	sess, err := NewTxListenerSession(ctx, l, chains, h)
	if err != nil {
		return err
	}

	// Wait for session exit signal
	<-sess.ctx.Done()

	return nil
}

// Close all listeners
func (l *TxListener) Close() {
	l.closeOnce.Do(func() {
		// Close listener
		log.Infof("tx-listener: closing...")

		close(l.closed)
	})
}

type TxListenerSession struct {
	ctx context.Context

	listener *TxListener

	h handler.TxListenerHandler

	chains []*big.Int
}

func NewTxListenerSession(ctx context.Context, l *TxListener, chains []*big.Int, h handler.TxListenerHandler) (*TxListenerSession, error) {
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	sess := &TxListenerSession{
		ctx:      cancelCtx,
		listener: l,
		h:        h,
		chains:   chains,
	}

	// Call handler Setup Hook
	err := h.Setup(sess)
	if err != nil {
		return nil, err
	}

	// Start listening each chain in separate goroutines
	wg := &sync.WaitGroup{}
	errors := make(chan error, len(chains))
	for _, chain := range chains {
		wg.Add(1)
		go func(chain *big.Int) {
			defer wg.Done()
			// Cancel session as soon as a first chain listener goroutine exits
			defer cancel()
			errors <- sess.listen(chain)
		}(chain)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errors:
		// Wait for go routines to complete
		wg.Wait()
		return sess, err
	case <-done:
		err = h.Cleanup(sess)
		if err != nil {
			return sess, err
		}
	}

	return sess, nil
}

func (s *TxListenerSession) Chains() []*big.Int {
	return s.chains
}

func (s *TxListenerSession) listen(chain *big.Int) error {
	select {
	case <-s.ctx.Done():
		return nil
	case <-s.listener.closed:
		return nil
	default:
	}

	tracker := tiptracker.NewTracker(s.listener.ec, chain, &s.listener.conf.TipTracker)

	// Retrieve initial position
	blockNumber, txIndex, err := s.h.GetInitialPosition(chain)
	if err != nil {
		return err
	}

	if blockNumber == -1 {
		blockNumber, _ = tracker.HighestBlock(context.Background())
		// Force tx index to 0
		txIndex = 0
	}

	// Create cursor
	cur := cursor.NewBlockCursorFromTracker(s.listener.ec, tracker, blockNumber, s.listener.conf.BlockCursor)

	listener := &ChainListener{
		listener:           s.listener,
		t:                  tracker,
		cur:                cur,
		initialBlockNumber: blockNumber,
		initialTxIndex:     txIndex,
		blocks:             make(chan *types.TxListenerBlock),
		receipts:           make(chan *types.TxListenerReceipt),
		errors:             make(chan *types.TxListenerError),
		blockNumber:        blockNumber,
		txIndex:            txIndex,
		ctx:                s.ctx,
		logger:             log.WithFields(log.Fields{"chain.id": chain.Text(10)}),
	}

	// Start cursor
	cur.Start()

	// Start listening
	go func() {
		_ = listener.listen()
	}()

	// Call Listen handler hook
	return s.h.Listen(s, listener)
}

// ChainListener listen to all transactions emitted from a chain
type ChainListener struct {
	listener *TxListener

	t   *tiptracker.Tracker
	cur *cursor.BlockCursor

	initialBlockNumber, initialTxIndex int64
	blockNumber, txIndex               int64

	blocks   chan *types.TxListenerBlock
	receipts chan *types.TxListenerReceipt
	errors   chan *types.TxListenerError

	ctx    context.Context
	logger *log.Entry
}

// ChainID returns Network ID of the chain being listened
func (l *ChainListener) ChainID() *big.Int {
	return l.t.ChainID()
}

// InitialPosition return initial position from which it started to listen
func (l *ChainListener) InitialPosition() (blockNumber, txIndex int64) {
	return l.initialBlockNumber, l.initialTxIndex
}

// Context return context associated to the listener session
func (l *ChainListener) Context() context.Context {
	return l.ctx
}

// Receipts returns a channel of Receipts as they are mined
func (l *ChainListener) Receipts() <-chan *types.TxListenerReceipt {
	return l.receipts
}

// Blocks returns a channel of Blocks as they are mined
func (l *ChainListener) Blocks() <-chan *types.TxListenerBlock {
	return l.blocks
}

// Errors returns a channel of Errors as they are mined
func (l *ChainListener) Errors() <-chan *types.TxListenerError {
	return l.errors
}

func (l *ChainListener) listen() error {
	// Start listener in a separate goroutine
	l.logger.WithFields(log.Fields{
		"start.block":    l.blockNumber,
		"start.tx-index": l.txIndex,
	}).Debug("listener: start loop")
feedingLoop:
	for {
		select {
		case <-l.ctx.Done():
			break feedingLoop
		case block, ok := <-l.cur.Blocks():
			if !ok {
				// Block cursor block channel has been closed so we leave the loop
				break feedingLoop
			}
			// We have a new block
			if l.listener.conf.TxListener.Return.Blocks {
				// This will be blocking until user consume from Blocks channel
				l.blocks <- block.Copy()
			}

			// We treat every transaction
			for l.txIndex < int64(len(block.Receipts)) {
				select {
				case <-l.ctx.Done():
					break feedingLoop
				default:
					l.receipts <- block.Receipts[l.txIndex]
					atomic.AddInt64(&l.txIndex, 1)
				}
			}

			// We have seen all receipts in current block so we prepare position for next block
			atomic.AddInt64(&l.blockNumber, 1)
			atomic.StoreInt64(&l.txIndex, 0)

		case err, ok := <-l.cur.Errors():
			if !ok {
				// Cursor error channel has been closed so we leave the loop
				break feedingLoop
			}

			// Send error
			if l.listener.conf.TxListener.Return.Errors {
				l.errors <- err
			} else {
				log.WithError(err).WithFields(log.Fields{
					"chain": l.t.ChainID().Text(10),
				}).Error("listener: failed to retrieve block")
			}
			return err
		}
	}
	l.logger.Debugf("listener: left loop")

	// Close channels
	close(l.receipts)
	close(l.blocks)
	close(l.errors)

	// Close cursor
	l.cur.Close()

	return nil
}
