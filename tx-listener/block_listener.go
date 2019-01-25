package listener

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
)

// BlockListener is an interface to get new block as they are mined
type BlockListener interface {
	// Blocks return a read channel of blocks
	Blocks() <-chan *types.Block

	// Errors return a read channel of errors
	Errors() <-chan error

	// Close stops consumer from fetching new blocks
	// It is required to call this function before a consumer object passes
	// out of scope, as it will otherwise leak memory.
	Close()
}

type blockListener struct {
	conf   *Config
	cur    BlockCursor
	blocks chan *types.Block
	errors chan error

	closeOnce       *sync.Once
	trigger, closed chan struct{}
}

func newBlockListener(cur BlockCursor, conf *Config) *blockListener {
	return &blockListener{
		conf:      conf,
		cur:       cur,
		blocks:    make(chan *types.Block),
		errors:    make(chan error),
		closeOnce: &sync.Once{},
		closed:    make(chan struct{}),
		trigger:   make(chan struct{}, 1),
	}
}

func (bl *blockListener) Blocks() <-chan *types.Block {
	return bl.blocks
}

func (bl *blockListener) Errors() <-chan error {
	return bl.errors
}

func (bl *blockListener) Close() {
	bl.closeOnce.Do(func() {
		close(bl.closed)
	})
}

func (bl *blockListener) feeder() {
	// Send a first message to trigger channel so feeder can start
	bl.trigger <- struct{}{}

	// Ticker allows to limit number of fetch calls on Ethereum client while waiting for a new block
	ticker := time.NewTicker(bl.conf.BlockListener.Backoff)
	defer ticker.Stop()

feedingLoop:
	for {
		select {
		case <-bl.closed:
			// Consumer is close thus we quit the loop
			break feedingLoop
		case <-bl.trigger:
			// Retrieve next mined block
			block, err := bl.cur.Next()
			if err != nil && bl.conf.BlockListener.Return.Errors {
				bl.errors <- err
			}

			if block != nil {
				// A new block has been mined
				bl.blocks <- block
				// A new mined block so we re-trigger in case next block has already been mined
				bl.trigger <- struct{}{}
			}
		case <-ticker.C:
			// Wait for a ticker then re-trigger
			if len(bl.trigger) > 0 {
				// If already triggered no need to trigger again
				continue feedingLoop
			}
			// Send an element to trigger channel
			bl.trigger <- struct{}{}
		}
	}
	close(bl.blocks)
	close(bl.errors)
	close(bl.trigger)
}
