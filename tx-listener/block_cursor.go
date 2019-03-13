package listener

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
)

// TxListenerReceipt contains useful information about a receipt
type TxListenerReceipt struct {
	// Chain receipt has been read from
	ChainID *big.Int

	// Go-Ethereum receipt
	types.Receipt

	// Position of the receipt
	BlockHash   common.Hash
	BlockNumber int64
	TxHash      common.Hash
	TxIndex     uint64
}

// TxListenerBlock contains data about a block
type TxListenerBlock struct {
	// Chain block has been read from
	ChainID *big.Int

	// Go-Ethereum block
	types.Block

	// Ordered receipts for every transaction in the block
	receipts []*TxListenerReceipt
}

// Copy creates a deep copy of a block to prevent side effects
func (b *TxListenerBlock) Copy() *TxListenerBlock {
	return &TxListenerBlock{
		ChainID:  big.NewInt(0).Set(b.ChainID),
		Block:    *b.WithBody(b.Transactions(), b.Uncles()),
		receipts: make([]*TxListenerReceipt, len(b.receipts)),
	}
}

// TxListenerError is what is provided to the user when an error occurs.
// It wraps an error and includes the chain ID
type TxListenerError struct {
	// Network ID the error occurred on
	ChainID *big.Int

	// Error
	Err error
}

func (e TxListenerError) Error() string {
	return fmt.Sprintf("tx-listener: error while listening on chain %s: %s", hexutil.EncodeBig(e.ChainID), e.Err)
}

// TxListenerErrors is a type that wraps a batch of errors and implements the Error interface.
type TxListenerErrors []*TxListenerErrors

func (e TxListenerErrors) Error() string {
	return fmt.Sprintf("tx-listener: %d errors while while listening", len(e))
}

// ChainTracker keep track of block chain highest mined block
type ChainTracker interface {
	ChainID() *big.Int
	HighestBlock(ctx context.Context) (int64, error)
}

// BaseTracker is a basic chain tracker
type BaseTracker struct {
	ec EthClient

	chainID *big.Int
	depth   uint64
}

// NewBaseTracker creates a new base tracker
func NewBaseTracker(ec EthClient, chainID *big.Int) *BaseTracker {
	return &BaseTracker{
		ec:      ec,
		chainID: chainID,
		depth:   uint64(viper.GetInt64("listener.tracker.depth")),
	}
}

// ChainID returns ID of the tracked chain
func (t *BaseTracker) ChainID() *big.Int {
	return big.NewInt(0).Set(t.chainID)
}

// HighestBlock returns highest mined block on the tracked chain
func (t *BaseTracker) HighestBlock(ctx context.Context) (int64, error) {
	header, err := t.ec.HeaderByNumber(ctx, t.chainID, nil)
	if err != nil {
		return 0, err
	}
	if header.Number.Uint64() <= t.depth {
		return 0, nil
	}
	return int64(header.Number.Uint64() - t.depth), nil
}

// Future is an element used to start a task and retrieve its result later
type Future struct {
	res chan interface{}
	err chan error
}

// Close future
func (f *Future) Close() {
	close(f.res)
	close(f.err)
}

// Cursor is an interface for a cursor object reading from a chain
type Cursor interface {
	// ChainID returns the chain ID the cursor is applied on
	ChainID() *big.Int

	// Current returns element the cursor is pointing on
	Blocks() <-chan *TxListenerBlock

	// Err returns a possible error met by the cursor when calling Next
	Errors() <-chan *TxListenerError

	// Close cursor
	Close()
}

// BlockCursor allows to retrieve new blocks
type BlockCursor struct {
	ec EthClient

	// Allows to get information about chain
	t ChainTracker

	// blockNumber stands last block that has been fetched
	// currentHead stands for the most advanced known mined block (we use it as cache so we do not need to always fetch Eth client for last mined block)
	blockNumber, currentHead int64

	// CurrentBlock on the cursor
	blocks chan *TxListenerBlock
	errors chan *TxListenerError

	// BlockCursor fetches blocks ahead of being consumed for performances matters
	// blockFutures is a channel of block being
	blockFutures chan *Future

	// Closing utilies
	closeOnce *sync.Once
	closed    chan struct{}
}

// NewBlockCursorFromTracker creates a new block cursor using a tracker
func NewBlockCursorFromTracker(ec EthClient, t ChainTracker, blockNumber int64) *BlockCursor {
	bc := newBlockCursorFromTracker(ec, t, blockNumber)

	// Start feeder & dispatcher in separate goroutines
	go bc.feeder()
	go bc.dispatcher()

	return bc
}

func newBlockCursorFromTracker(ec EthClient, t ChainTracker, blockNumber int64) *BlockCursor {
	return &BlockCursor{
		ec:           ec,
		t:            t,
		blockNumber:  blockNumber,
		currentHead:  0,
		blocks:       make(chan *TxListenerBlock),
		errors:       make(chan *TxListenerError),
		blockFutures: make(chan *Future, uint64(viper.GetInt64("listener.block.limit"))),
		closed:       make(chan struct{}),
		closeOnce:    &sync.Once{},
	}
}

// NewBlockCursor creates a new block cursor for a given chain starting at a given blockNumber
func NewBlockCursor(ec EthClient, chainID *big.Int, blockNumber int64) *BlockCursor {
	return NewBlockCursorFromTracker(ec, NewBaseTracker(ec, chainID), blockNumber)
}

// ChainID returns ID of the chain the cursor is applied on
func (bc *BlockCursor) ChainID() *big.Int {
	return bc.t.ChainID()
}

// Next moves cursor to Next available block
func (bc *BlockCursor) dispatcher() {
	// In case a future block is available we treat it
	for future := range bc.blockFutures {
		select {
		case err := <-future.err:
			bc.errors <- err.(*TxListenerError)
		case res := <-future.res:
			bc.blocks <- res.(*TxListenerBlock)
		}
	}
	close(bc.blocks)
	close(bc.errors)
}

// Blocks return channel of blocks
func (bc *BlockCursor) Blocks() <-chan *TxListenerBlock {
	return bc.blocks
}

// Errors return channel of errors
func (bc *BlockCursor) Errors() <-chan *TxListenerError {
	return bc.errors
}

// Close cursor
func (bc *BlockCursor) Close() {
	bc.closeOnce.Do(func() {
		close(bc.closed)
	})
}

// feeder feeds the blockFutures channel
// It manages the main cursor loop that fetch blocks & receipts from Eth client
func (bc *BlockCursor) feeder() {
	ctx, cancel := context.WithCancel(context.Background())
feedingLoop:
	for {
		select {
		case <-bc.closed:
			// Cancel pending fetches and break loop
			cancel()
			break feedingLoop
		default:
			if bc.blockNumber <= bc.currentHead {
				// We are behind chain head meaning we have at leasdt one block to fetch
				bc.blockFutures <- bc.fetchBlock(ctx, bc.blockNumber)
				bc.blockNumber++
			} else {
				// We are ahead of last known chain head, so we refresh it
				head, err := bc.t.HighestBlock(ctx)
				if head > bc.currentHead {
					// Chain has moved forward (meaning new blocks have been mined and are ready to be fetched)
					bc.currentHead = head
				} else {
					// We are still ahead or something went wrong
					if err != nil {
						// If we got an error while retrieving chain head we notify it in a future
						bFuture := &Future{
							res: make(chan interface{}),
							err: make(chan error),
						}

						go func(err error) {
							// Notify error and Close future
							defer bFuture.Close()
							bFuture.err <- err
						}(err)

						bc.blockFutures <- bFuture
					}

					// Chain has not move forward so we sleep before retrying (waiting for updates on the chain)
					time.Sleep(viper.GetDuration("listener.block.backoff"))
				}
			}
		}
	}
	close(bc.blockFutures)
}

// fetchBlock fetch a block asynchronously
func (bc *BlockCursor) fetchBlock(ctx context.Context, blockNumber int64) *Future {
	bFuture := &Future{
		res: make(chan interface{}),
		err: make(chan error),
	}

	// Retrieve block in a separate goroutine
	go func() {
		defer bFuture.Close()

		block, err := bc.ec.BlockByNumber(ctx, bc.ChainID(), big.NewInt(blockNumber))
		if err != nil {
			bFuture.err <- err
			return
		}

		// Block should be available
		if block == nil {
			bFuture.err <- BlockMissingError(blockNumber)
			return
		}

		bl := &TxListenerBlock{
			ChainID:  bc.ChainID(),
			Block:    *block,
			receipts: []*TxListenerReceipt{},
		}

		// Retrieve receipts in separate go-routines
		rFutures := []*Future{}
		for _, tx := range bl.Block.Transactions() {
			rFutures = append(rFutures, bc.fetchReceipt(ctx, tx.Hash()))
		}

		for i, rFuture := range rFutures {
			select {
			case err := <-rFuture.err:
				bFuture.err <- err
				return
			case res := <-rFuture.res:
				receipt := res.(*TxListenerReceipt)
				receipt.TxIndex = uint64(i)
				receipt.BlockHash = block.Hash()
				receipt.BlockNumber = block.Number().Int64()
				bl.receipts = append(bl.receipts, receipt)
			}
		}

		// Return block in result
		bFuture.res <- bl
	}()

	return bFuture
}

// fetchReceipt fetch a receipt asynchronously
func (bc *BlockCursor) fetchReceipt(ctx context.Context, txHash common.Hash) *Future {
	future := &Future{
		res: make(chan interface{}),
		err: make(chan error),
	}

	go func() {
		defer future.Close()
		receipt, err := bc.ec.TransactionReceipt(ctx, bc.ChainID(), txHash)
		if err != nil {
			future.err <- err
			return
		}

		if receipt == nil {
			future.err <- ReceiptMissingError(txHash.Hex())
			return
		}

		r := &TxListenerReceipt{
			ChainID: bc.ChainID(),
			Receipt: *receipt,
			TxHash:  txHash,
		}
		future.res <- r
	}()

	return future
}
