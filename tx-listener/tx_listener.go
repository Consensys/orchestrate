package listener

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

// TxListener is main inferface
type TxListener interface {
	// Listen creates a ChainListener on the the given CHain network
	// It will return an error if this TxListener is already listening
	// on the given chain
	// BlockNumber can be a blockNumber or -1 for latest block
	Listen(chainID *big.Int, blockNumber int64, txIndex int64, conf Config) (ChainListener, error)

	// Receipts returns a read channel of receipts that are returned by the TxListener
	Receipts() <-chan *TxListenerReceipt

	// Receipts returns a read channel of blocks that are returned by the TxListener
	Blocks() <-chan *TxListenerBlock

	// Receipts returns a read channel of errors that are returned by the TxListener
	Errors() <-chan *TxListenerError

	// Chains return a list of chains that are currently listen
	Chains() []*big.Int

	// Progress returns progress of TxListener
	Progress(ctx context.Context) map[string]*Progress

	// Close TxListener
	Close()
}

// NewTxListener creates a new TxListener
func NewTxListener(ec EthClient) TxListener {
	return &txListener{
		ec:             ec,
		mux:            &sync.RWMutex{},
		chainListeners: make(map[string]*singleChainListener),
		blocks:         make(chan *TxListenerBlock),
		receipts:       make(chan *TxListenerReceipt),
		errors:         make(chan *TxListenerError),
		wait:           &sync.WaitGroup{},
		closeOnce:      &sync.Once{},
		closed:         false,
	}
}

type txListener struct {
	ec EthClient

	mux            *sync.RWMutex
	chainListeners map[string]*singleChainListener

	blocks   chan *TxListenerBlock
	receipts chan *TxListenerReceipt
	errors   chan *TxListenerError

	wait      *sync.WaitGroup
	closeOnce *sync.Once
	closed    bool
}

func (l *txListener) Listen(chainID *big.Int, blockNumber int64, txIndex int64, conf Config) (ChainListener, error) {
	t := &BaseTracker{
		ec:      l.ec,
		chainID: chainID,
		depth:   conf.BlockCursor.Tracker.Depth,
	}

	cur := newBlockCursorFromTracker(l.ec, t, blockNumber, conf)
	listener := &singleChainListener{
		t:           t,
		cur:         cur,
		conf:        conf,
		blocks:      l.blocks,
		receipts:    l.receipts,
		errors:      l.errors,
		blockNumber: blockNumber,
		txIndex:     txIndex,
		closeOnce:   &sync.Once{},
		closed:      make(chan struct{}),
	}

	// Register new listener
	err := l.addListener(listener)
	if err != nil {
		listener.Close()
		return nil, err
	}

	// Start listener in a separate goroutine
	log.WithFields(log.Fields{
		"Chain": chainID.Text(16),
	}).Infof("Start listening from blockNumber=%v txIndex=%v", blockNumber, txIndex)

	l.wait.Add(1)
	go listener.feeder()
	go cur.feeder()

	return listener, nil
}

// Receipt return a channel of receipts
func (l *txListener) Receipts() <-chan *TxListenerReceipt {
	return l.receipts
}

// Receipt return a channel of receipts
func (l *txListener) Blocks() <-chan *TxListenerBlock {
	return l.blocks
}

// Receipt return a channel of receipts
func (l *txListener) Errors() <-chan *TxListenerError {
	return l.errors
}

// Chains returns Network IDs that are currently listened
func (l *txListener) Chains() []*big.Int {
	l.mux.RLock()
	defer l.mux.RUnlock()
	chains := []*big.Int{}
	for _, listener := range l.chainListeners {
		chains = append(chains, listener.ChainID())
	}
	return chains
}

// Progress return progress for every chains
func (l *txListener) Progress(ctx context.Context) map[string]*Progress {
	progress := make(map[string]*Progress)
	l.mux.RLock()
	defer l.mux.RUnlock()
	for chainID, listener := range l.chainListeners {
		progress[chainID], _ = listener.Progress(ctx)
	}

	return progress
}

// Close all listeners
func (l *txListener) Close() {
	l.closeOnce.Do(func() {
		// Close listener
		log.Infof("Closing TxListener...")

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
		log.Infof("TxListener closed")
	})
}

func (l *txListener) addListener(listener *singleChainListener) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.closed {
		return fmt.Errorf("Listener has been closed")
	}

	chainID := listener.ChainID().Text(16)
	_, ok := l.chainListeners[chainID]
	if ok {
		return fmt.Errorf("Chain %q is already being listened", chainID)
	}
	listener.txlistener = l
	l.chainListeners[chainID] = listener

	return nil
}

func (l *txListener) removeListener(listener *singleChainListener) {
	l.mux.Lock()
	defer l.mux.Unlock()

	chainID := listener.ChainID().Text(16)
	delete(l.chainListeners, chainID)
}

// Progress holds information about listener progress
type Progress struct {
	CurrentBlock int64 // Current block number where the listener is
	TxIndex      int64 // Current txIndex where the listener is
	HighestBlock int64 // Highest alleged block number in the chain
}

// ChainListener is a listener that listens for a given chain
type ChainListener interface {
	ChainID() *big.Int
	Progress(ctx context.Context) (*Progress, error)
	Close()
}

// singleChainListener listen to all transactions emitted from a chain
type singleChainListener struct {
	txlistener *txListener
	t          ChainTracker
	cur        Cursor

	conf Config

	mux                  *sync.RWMutex
	blockNumber, txIndex int64

	blocks   chan<- *TxListenerBlock
	receipts chan<- *TxListenerReceipt
	errors   chan<- *TxListenerError

	// Closing utilies
	closeOnce *sync.Once
	closed    chan struct{}
}

// ChainID returns Network ID of the chain being listened
func (l *singleChainListener) ChainID() *big.Int {
	return l.t.ChainID()
}

// Progress returns current listener progress
func (l *singleChainListener) Progress(ctx context.Context) (*Progress, error) {
	head, err := l.t.HighestBlock(ctx)
	if err != nil {
		return &Progress{atomic.LoadInt64(&l.blockNumber), atomic.LoadInt64(&l.txIndex), -1}, err
	}
	return &Progress{atomic.LoadInt64(&l.blockNumber), atomic.LoadInt64(&l.txIndex), head}, nil
}

func (l *singleChainListener) Close() {
	l.closeOnce.Do(func() {
		// Close listener
		log.WithFields(log.Fields{
			"Chain": l.t.ChainID().Text(16),
		}).Infof("Closing listener...")
		close(l.closed)
	})
}

func (l *singleChainListener) feeder() {
	log.WithFields(log.Fields{
		"Chain": l.t.ChainID().Text(16),
	}).Debugf("tx-listener: start listening loop")

feedingLoop:
	for {
		select {
		case <-l.closed:
			break feedingLoop
		default:
			if l.cur.Current() != nil && l.txIndex < int64(len(l.cur.Current().receipts)) {
				l.receipts <- l.cur.Current().receipts[l.txIndex]
				atomic.AddInt64(&l.txIndex, 1)
				if l.txIndex == int64(len(l.cur.Current().receipts)) {
					// We have seen all receipts in current block so we prepare for next block
					atomic.AddInt64(&l.blockNumber, 1)
					atomic.StoreInt64(&l.txIndex, 0)
				} else {
					continue
				}
			}

			// Try to retrieve next block
			ok := l.cur.Next(context.Background())
			if !ok {
				// No new block available
				// Do we have an error?
				if err := l.cur.Err(); err != nil {
					// Send error
					if l.conf.TxListener.Return.Errors {
						l.errors <- err
					} else {
						log.WithFields(log.Fields{
							"Chain": l.t.ChainID().Text(16),
						}).Error(err.Error())
					}
					// We abort the listener
					l.Close()
					break feedingLoop
				}
				continue
			}

			// We have a new block
			if l.conf.TxListener.Return.Blocks {
				l.blocks <- l.cur.Current().Copy()
			} else {
				log.WithFields(log.Fields{
					"Chain": l.t.ChainID().Text(16),
				}).Debugf("tx-listener: New block %v", l.cur.Current().Hash().Hex())
			}
		}
	}

	log.WithFields(log.Fields{
		"Chain": l.t.ChainID().Text(16),
	}).Debugf("tx-listener: left listening loop")

	// Close cursor if not nil
	if l.cur != nil {
		l.cur.Close()
	}

	// We indicate that feeder has stoped
	if l.txlistener != nil {
		l.txlistener.removeListener(l)
	}

	log.WithFields(log.Fields{
		"Chain": l.t.ChainID().Text(16),
	}).Infof("Listener closed...")

	l.txlistener.wait.Done()
}
