package base

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/ethclient"
	cursor "gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/block-cursor/base"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/handler"
	tiptracker "gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/tx-listener/tip-tracker/base"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
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
	sess := NewTxListenerSession(ctx, l, chains, h)

	// Start session
	return sess.run()
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

	h      handler.TxListenerHandler
	cancel func()
	chains []*big.Int
}

func NewTxListenerSession(ctx context.Context, l *TxListener, chains []*big.Int, h handler.TxListenerHandler) *TxListenerSession {
	log.Info("tx-listener: creating new listener session: ", chains)
	cancelCtx, cancel := context.WithCancel(ctx)

	return &TxListenerSession{
		ctx:      cancelCtx,
		cancel:   cancel,
		listener: l,
		h:        h,
		chains:   chains,
	}
}

func (s *TxListenerSession) run() error {
	// Call handler Setup Hook
	err := s.h.Setup(s)
	if err != nil {
		return err
	}

	// Start listening each chain in separate goroutines
	wg := &sync.WaitGroup{}
	errors := make(chan error, len(s.chains))
	defer close(errors)
	for _, chain := range s.chains {
		wg.Add(1)
		go func(chain *big.Int) {
			defer wg.Done()
			// Cancel session as soon as a first chain listener goroutine exits
			defer s.cancel()
			errors <- s.listen(chain)
		}(chain)
	}

	// Wait for all listening go routines to complete
	wg.Wait()

	// Clean
	err = s.h.Cleanup(s)
	if err != nil {
		return err
	}

	select {
	case e := <-errors:
		err = e
	default:
	}

	return err
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

	log.Infof("tx-listener: getting initial position in chain %s", chain.String())
	// Retrieve initial position
	blockNumber, txIndex, err := s.h.GetInitialPosition(chain)
	if err != nil {
		log.WithError(err).Errorf("tx-listener: failed to get initial position in chain %s", chain.String())
		return err
	}

	if blockNumber == -1 {
		log.Infof("tx-listener: getting highest block number")
		blockNumber, err = tracker.HighestBlock(context.Background())
		if err != nil {
			log.WithError(err).Errorf("tx-listener: failed to get highest block number")
			return err
		}
		// Force tx index to 0
		txIndex = 0
	}

	// Create cursor
	log.Infof("tx-listener: creating cursor from tracker")
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
	}).Debug("tx-listener: start loop")
feedingLoop:
	for {
		select {
		case <-l.ctx.Done():
			log.Warnf("tx-listener: loop for chain %s is over", l.ChainID().String())
			break feedingLoop
		case block, ok := <-l.cur.Blocks():
			if !ok {
				log.Warnf("tx-listener: Block cursor block channel has been closed. Closing blocks listener loop for chain: %s", l.ChainID().String())
				break feedingLoop
			}
			// We have a new block
			if l.listener.conf.TxListener.Return.Blocks {
				log.Debugf("tx-listener: get a new block for chain: %s", l.ChainID().String())
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
				log.Warnf("tx-listener: block cursor block channel has been closed. Closing blocks listener loop for chain: %s", l.ChainID().String())
				break feedingLoop
			}

			// Send error
			if l.listener.conf.TxListener.Return.Errors {
				l.errors <- err
			} else {
				log.WithError(err).WithFields(log.Fields{
					"chain": l.t.ChainID().Text(10),
				}).Error("tx-listener: failed to retrieve block")
			}
			return err
		}
	}
	l.logger.Debug("tx-listener: left loop")

	// Close channels
	close(l.receipts)
	close(l.blocks)
	close(l.errors)

	// Close cursor
	l.cur.Close()

	return nil
}
