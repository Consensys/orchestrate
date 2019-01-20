package ethereum

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// BlockCursor allows to retrieve a new block
type BlockCursor interface {
	Next() (*types.Block, error)
	Set(pos *big.Int)
}

// TxListenerClient is a minimal EthClient interface required by a TxListener
type TxListenerClient interface {
	// BlockByNumber retrieve a block by its number
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)

	// TransactionReceipt retrieve a transaction receipt given a transaction hash
	TransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error)
}

type blockCursor struct {
	c TxListenerClient

	next *big.Int
	// TODO: add an history of blocks of configurable length
	// (we could add some checks to ensure no re-org happened)
}

func newBlockCursor(c TxListenerClient) *blockCursor {
	cursor := &blockCursor{
		c:    c,
		next: big.NewInt(0),
	}
	return cursor
}

// Next returns next block mined if available
func (bc *blockCursor) Next() (*types.Block, error) {
	// Retrieve next block
	block, err := bc.c.BlockByNumber(context.Background(), bc.next)
	if err != nil {
		return nil, err
	}

	// If we retrieved a block we increment cursor position
	if block != nil {
		bc.next = bc.next.Add(bc.next, big.NewInt(1))
	}

	return block, nil
}

// Set position of cursor
func (bc *blockCursor) Set(pos *big.Int) {
	bc.next.Set(pos)
}

// BlockConsumer is an interface to get new block as they are mined
type BlockConsumer interface {
	// Blocks return a read channel of blocks
	Blocks() <-chan *types.Block

	// Close stops consumer from fetching new blocks
	// It is required to call this function before a consumer object passes
	// out of scope, as it will otherwise leak memory.
	Close()
}

type blockConsumer struct {
	conf   *TxListenerConfig
	cur    BlockCursor
	blocks chan *types.Block

	closeOnce       *sync.Once
	trigger, closed chan struct{}
}

func newBlockConsumer(cur BlockCursor, conf *TxListenerConfig) *blockConsumer {
	return &blockConsumer{
		conf:      conf,
		cur:       cur,
		blocks:    make(chan *types.Block),
		closeOnce: &sync.Once{},
		closed:    make(chan struct{}),
		trigger:   make(chan struct{}, 1),
	}
}

func (c *blockConsumer) Blocks() <-chan *types.Block {
	return c.blocks
}

func (c *blockConsumer) Close() {
	c.closeOnce.Do(func() {
		close(c.closed)
	})
}

func (c *blockConsumer) feeder() {
	// Ticker allows to limit number of fetch calls on Ethereum client while waiting for a new block
	ticker := time.NewTicker(c.conf.Block.Delay)
	defer ticker.Stop()

	// Trigger execution
	c.trigger <- struct{}{}
feedingLoop:
	for {
		select {
		case <-c.closed:
			// Consumer is close thus we quit the loop
			break feedingLoop
		case <-c.trigger:
			block, err := c.cur.Next()
			if err != nil {
				// TODO: implement backoff retry strategy
				c.Close()
			}
			if block != nil {
				// A new block to listen
				c.blocks <- block
				// We got a new block so we re-trigger in case next block has already been mined
				c.trigger <- struct{}{}
			}
		case <-ticker.C:
			if len(c.trigger) > 0 {
				continue
			}
			// We re-trigger execution
			c.trigger <- struct{}{}
		}
	}
	close(c.blocks)
	close(c.trigger)
}

// TxListenerConfig configuration of a TxListener
type TxListenerConfig struct {
	Block struct {
		// Delay to wait between calls to get new mined blocks
		Delay time.Duration
	}

	Receipts struct {
		// Count indicate how many receipts can be retrieved in parallel
		Count uint
	}
}

// TxListener is an interface to listen to a Ethereum blockchain activity
type TxListener interface {
	// Blocks return a read channel of blocks
	Blocks() <-chan *types.Block

	// Receipts return a read channel of transaction receipts
	Receipts() <-chan *types.Receipt

	// Close stops listener
	Close()
}

type receiptResponse struct {
	err chan error
	res chan *types.Receipt
}

type txListener struct {
	ec TxListenerClient

	receipts chan *types.Receipt
	blocks   chan *types.Block

	bc        BlockConsumer
	responses chan *receiptResponse

	closeOnce *sync.Once
	closed    chan struct{}
}

// NewTxListener creates a new transaction listener
func NewTxListener(ec TxListenerClient, conf *TxListenerConfig) TxListener {
	// Instantiate txListener
	cur := newBlockCursor(ec)
	bc := newBlockConsumer(cur, conf)
	l := &txListener{
		ec:        ec,
		receipts:  make(chan *types.Receipt),
		blocks:    make(chan *types.Block),
		bc:        bc,
		responses: make(chan *receiptResponse, conf.Receipts.Count),
		closeOnce: &sync.Once{},
		closed:    make(chan struct{}),
	}

	// Start feeding txListener
	go l.feeder()
	go l.dispatcher()

	// Start feeding block consumer
	go bc.feeder()

	return l
}

func (l *txListener) Blocks() <-chan *types.Block {
	return l.blocks
}

func (l *txListener) Receipts() <-chan *types.Receipt {
	return l.receipts
}

func (l *txListener) dispatcher() {
dispatchLoop:
	for block := range l.bc.Blocks() {
		// A new block has been mined
		// Send block to blocks channel
		l.blocks <- block

		// Retrieve transaction receipt for every transactions
		for _, tx := range block.Transactions() {
			select {
			case <-l.closed:
				break dispatchLoop
			default:
				// We retrieve receipt concurrently
				l.responses <- l.getReceipt(tx)
			}
		}
	}
	close(l.responses)
	close(l.blocks)
}

func (l *txListener) feeder() {
	for response := range l.responses {
		// Retrieve receipt response and wait for it to complete
		select {
		case <-response.err:
			// TODO: handle error case
			l.Close()
		case receipt := <-response.res:
			l.receipts <- receipt
		}
	}
	close(l.receipts)
}

func (l *txListener) getReceipt(tx *types.Transaction) *receiptResponse {
	response := &receiptResponse{
		err: make(chan error),
		res: make(chan *types.Receipt),
	}
	// Retrieve receipt in a parallel goroutine
	go func() {
		receipt, err := l.ec.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			response.err <- err
		} else {
			response.res <- receipt
		}
	}()
	return response
}

func (l *txListener) Close() {
	l.closeOnce.Do(func() {
		close(l.closed)
		l.bc.Close()
	})
}
