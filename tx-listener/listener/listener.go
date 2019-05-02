package listener

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	cursor "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener/block-cursor/base"
	tiptracker "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener/tip-tracker/base"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/types"
)

type TxListener interface {
	// Listen creates a ChainListener on the the given CHain network
	// It will return an error if this TxListener is already listening
	// on the given chain
	// BlockNumber can be a blockNumber or -1 for latest block
	Listen(chainID *big.Int, blockNumber int64, txIndex int64) (ChainListener, error)

	// Receipts returns a read channel of receipts that are returned by the TxListener
	Receipts() <-chan *types.TxListenerReceipt

	// Receipts returns a read channel of blocks that are returned by the TxListener
	Blocks() <-chan *types.TxListenerBlock

	// Receipts returns a read channel of errors that are returned by the TxListener
	Errors() <-chan *types.TxListenerError

	// Chains return a list of chains that are currently listen
	Chains() []*big.Int

	// Progress returns progress of TxListener
	Progress(ctx context.Context) map[string]*types.Progress

	// Close TxListener
	Close()
}

type listener struct {
	ec ethclient.ChainLedgerReader

	mux            *sync.RWMutex
	chainListeners map[string]*singleChainListener

	blocks   chan *types.TxListenerBlock
	receipts chan *types.TxListenerReceipt
	errors   chan *types.TxListenerError

	wait      *sync.WaitGroup
	closeOnce *sync.Once
	closed    bool

	conf *Config
}

// NewListener creates a new Listener
func NewListener(ec ethclient.ChainLedgerReader, conf *Config) TxListener {
	return &listener{
		ec:             ec,
		mux:            &sync.RWMutex{},
		chainListeners: make(map[string]*singleChainListener),
		blocks:         make(chan *types.TxListenerBlock),
		receipts:       make(chan *types.TxListenerReceipt),
		errors:         make(chan *types.TxListenerError),
		wait:           &sync.WaitGroup{},
		closeOnce:      &sync.Once{},
		closed:         false,
		conf:           conf,
	}
}

func (l *listener) Listen(chainID *big.Int, blockNumber, txIndex int64) (ChainListener, error) {
	// Set chain tracker
	tracker := tiptracker.NewTracker(l.ec, chainID, &l.conf.TipTracker)

	// Set cursor
	if blockNumber == -1 {
		// We start from highest block
		blockNumber, _ = tracker.HighestBlock(context.Background())
	}
	cur := cursor.NewBlockCursorFromTracker(l.ec, tracker, blockNumber, l.conf.BlockCursor)

	// Create listener
	listener := &singleChainListener{
		t:           tracker,
		cur:         cur,
		blocks:      l.blocks,
		receipts:    l.receipts,
		errors:      l.errors,
		blockNumber: blockNumber,
		txIndex:     txIndex,
		closeOnce:   &sync.Once{},
		closed:      make(chan struct{}),
		conf:        l.conf,
	}

	// Register new listener
	err := l.addListener(listener)
	if err != nil {
		close(listener.closed)
		return nil, err
	}

	// Start feeders in separate go routine
	l.wait.Add(1)
	go listener.feeder()
	cur.Start()

	return listener, nil
}

// Receipt return a channel of receipts
func (l *listener) Receipts() <-chan *types.TxListenerReceipt {
	return l.receipts
}

// Receipt return a channel of receipts
func (l *listener) Blocks() <-chan *types.TxListenerBlock {
	return l.blocks
}

// Receipt return a channel of receipts
func (l *listener) Errors() <-chan *types.TxListenerError {
	return l.errors
}

// Chains returns Network IDs that are currently listened
func (l *listener) Chains() []*big.Int {
	l.mux.RLock()
	defer l.mux.RUnlock()
	chains := []*big.Int{}
	for _, listener := range l.chainListeners {
		chains = append(chains, listener.ChainID())
	}
	return chains
}

// Progress return progress for every chains
func (l *listener) Progress(ctx context.Context) map[string]*types.Progress {
	progress := make(map[string]*types.Progress)
	l.mux.RLock()
	defer l.mux.RUnlock()
	for chainID, listener := range l.chainListeners {
		progress[chainID], _ = listener.Progress(ctx)
	}

	return progress
}

// Close all listeners
func (l *listener) Close() {
	l.closeOnce.Do(func() {
		// Close listener
		log.Infof("tx-listener: closing...")

		// Close every channel
		l.mux.Lock()
		for _, listener := range l.chainListeners {
			listener.Close()
		}
		l.closed = true
		l.mux.Unlock()

		// Wait for listeners to complete then close channels
		l.wait.Wait()
		close(l.receipts)
		close(l.blocks)
		close(l.errors)

		// Close listener
		log.Infof("tx-listener: closed")
	})
}

func (l *listener) addListener(listener *singleChainListener) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.closed {
		return fmt.Errorf("listener has been closed")
	}

	chainID := listener.ChainID().Text(16)
	_, ok := l.chainListeners[chainID]
	if ok {
		return fmt.Errorf("chain %q is already being listened", chainID)
	}
	listener.txlistener = l
	l.chainListeners[chainID] = listener

	return nil
}

func (l *listener) removeListener(listener *singleChainListener) {
	l.mux.Lock()
	defer l.mux.Unlock()

	chainID := listener.ChainID().Text(16)
	delete(l.chainListeners, chainID)

	// If no listener remaining then close (in parallel go routine to avoid dead lock)
	if len(l.chainListeners) == 0 {
		go l.Close()
	}
}

// ChainListener is a listener that listens for a given chain
type ChainListener interface {
	ChainID() *big.Int
	Progress(ctx context.Context) (*types.Progress, error)
	Close()
}

// singleChainListener listen to all transactions emitted from a chain
type singleChainListener struct {
	txlistener *listener
	t          *tiptracker.Tracker
	cur        *cursor.BlockCursor

	conf *Config

	blockNumber, txIndex int64

	blocks   chan<- *types.TxListenerBlock
	receipts chan<- *types.TxListenerReceipt
	errors   chan<- *types.TxListenerError

	// Closing utilies
	closeOnce *sync.Once
	closed    chan struct{}
}

// ChainID returns Network ID of the chain being listened
func (l *singleChainListener) ChainID() *big.Int {
	return l.t.ChainID()
}

// Progress returns current listener progress
func (l *singleChainListener) Progress(ctx context.Context) (*types.Progress, error) {
	head, err := l.t.HighestBlock(ctx)
	if err != nil {
		return &types.Progress{
			CurrentBlock: atomic.LoadInt64(&l.blockNumber),
			TxIndex:      atomic.LoadInt64(&l.txIndex),
			HighestBlock: -1,
		}, err
	}
	return &types.Progress{
		CurrentBlock: atomic.LoadInt64(&l.blockNumber),
		TxIndex:      atomic.LoadInt64(&l.txIndex),
		HighestBlock: head,
	}, nil
}

func (l *singleChainListener) Close() {
	l.closeOnce.Do(func() {
		// Close listener
		log.WithFields(log.Fields{
			"Chain": l.t.ChainID().Text(16),
		}).Infof("tx-listener: stop listening chain...")
		close(l.closed)
	})
}

func (l *singleChainListener) feeder() {
	// Start listener in a separate goroutine
	log.WithFields(log.Fields{
		"Chain": l.ChainID().Text(16),
	}).Infof("tx-listener: start listening from block=%v tx=%v", l.blockNumber, l.txIndex)
feedingLoop:
	for {
		select {
		case <-l.closed:
			break feedingLoop
		case block, ok := <-l.cur.Blocks():
			if !ok {
				// Block cursor block channel has been closed so we leave the loop
				break feedingLoop
			}
			// We have a new block
			if l.conf.TxListener.Return.Blocks {
				// This will be blocking until user consume from Blocks channel
				l.blocks <- block.Copy()
			}

			// We treat every transaction
			for l.txIndex < int64(len(block.Receipts)) {
				select {
				case <-l.closed:
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
				// Block cursor block channel has been closed so we leave the loop
				break feedingLoop
			}

			// Send error
			if l.conf.TxListener.Return.Errors {
				l.errors <- err
			} else {
				log.WithError(err).WithFields(log.Fields{
					"Chain": l.t.ChainID().Text(16),
				}).Error("Failed to retrieve block")
			}

			// We got an error so we abort the listener
			l.Close()
			break feedingLoop
		}
	}

	// Close cursor
	l.cur.Close()

	// We notify main txlistener that we closed
	if l.txlistener != nil {
		l.txlistener.removeListener(l)
	}

	l.txlistener.wait.Done()
	l.txlistener = nil

	log.WithFields(log.Fields{
		"Chain": l.t.ChainID().Text(16),
	}).Infof("tx-listener: closed")
}
