package base

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/logger"
	tiptracker "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/tx-listener/tip-tracker"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
)

// BlockCursor allows to retrieve new blocks
type BlockCursor struct {
	ec ethclient.ChainLedgerReader

	// Allows to get information about chain
	t tiptracker.TipTracker

	// blockNumber stands last block that has been fetched
	// currentHead stands for the most advanced known mined block (we use it as cache so we do not need to always fetch Eth client for last mined block)
	blockNumber, currentHead int64

	// ticker and trigger are used to control the flow of fetch call for new mined blocks
	ticker  *time.Ticker
	trigger chan struct{}

	// CurrentBlock on the cursor
	blocks chan *types.TxListenerBlock
	errors chan *types.TxListenerError

	// BlockCursor fetches blocks ahead of being consumed for performances matters
	// blockFutures is a channel of block being
	blockFutures chan *types.Future

	// Closing utils
	closeOnce *sync.Once
	closed    chan struct{}

	conf Config
}

func NewBlockCursorFromTracker(ec ethclient.ChainLedgerReader, t tiptracker.TipTracker, startBlockNumber int64, conf Config) *BlockCursor {
	return &BlockCursor{
		ec:           ec,
		t:            t,
		blockNumber:  startBlockNumber,
		currentHead:  0,
		blocks:       make(chan *types.TxListenerBlock),
		errors:       make(chan *types.TxListenerError),
		ticker:       time.NewTicker(conf.Backoff),
		trigger:      make(chan struct{}, 1),
		blockFutures: make(chan *types.Future, conf.Limit),
		closed:       make(chan struct{}),
		closeOnce:    &sync.Once{},
		conf:         conf,
	}
}

// Start start the cursor
func (bc *BlockCursor) Start() {
	go bc.feeder()
	go bc.dispatcher()
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
		case err := <-future.Err():
			bc.errors <- err.(*types.TxListenerError)
		case res := <-future.Result():
			bc.blocks <- res.(*types.TxListenerBlock)
		}
	}
	close(bc.blocks)
	close(bc.errors)
}

// Blocks return channel of blocks
func (bc *BlockCursor) Blocks() <-chan *types.TxListenerBlock {
	return bc.blocks
}

// Errors return channel of errors
func (bc *BlockCursor) Errors() <-chan *types.TxListenerError {
	return bc.errors
}

// Close cursor
func (bc *BlockCursor) Close() {
	bc.closeOnce.Do(func() {
		close(bc.closed)
	})
}

func (bc *BlockCursor) trig() {
	select {
	case bc.trigger <- struct{}{}:
	default:
		// already triggered
	}
}

// feeder feeds the blockFutures channel
// It manages the main cursor loop that fetch blocks & receipts from Eth client
func (bc *BlockCursor) feeder() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// We trigger to start listener
	bc.trig()
feedingLoop:
	for {
		select {
		case <-bc.closed:
			// Cancel pending fetches and break loop
			cancel()
			break feedingLoop
		case <-bc.trigger:
			if bc.blockNumber <= bc.currentHead {
				// We are behind chain head meaning we have at least one block to fetch
				bc.blockFutures <- bc.fetchBlock(ctx, bc.blockNumber)

				// Increment block position
				bc.blockNumber++

				// As a block was available we re-trigger
				bc.trig()
			} else {
				// We are ahead of chain head so we refresh chain head
				head, err := bc.t.HighestBlock(ctx)
				log.WithFields(log.Fields{
					"number": head,
				}).Tracef("tx-listener: highest block")
				if head > bc.currentHead {
					// Chain has moved forward (meaning at least one new final block is ready to be fetched)
					bc.currentHead = head

					// We trigger
					bc.trig()
				} else if err != nil {
					// If we got an error while retrieving chain highest final block we notify it in a future
					bFuture := types.NewFuture()

					go func(err error) {
						// Notify error and Close future
						defer bFuture.Close()
						bFuture.Err() <- err
					}(err)

					bc.blockFutures <- bFuture
				}
			}
		case <-bc.ticker.C:
			// We trigger on every tick
			bc.trig()
		}
	}
	close(bc.blockFutures)
	close(bc.trigger)
	bc.ticker.Stop()
}

// fetchBlock fetch a block asynchronously
func (bc *BlockCursor) fetchBlock(ctx context.Context, blockNumber int64) *types.Future {
	bFuture := types.NewFuture()

	log.WithFields(log.Fields{
		"block.number": blockNumber,
	}).Tracef("tx-listener: fetch block")
	// Retrieve block in a separate goroutine
	go func() {
		defer bFuture.Close()
		logCtx := logger.WithLogEntry(
			ctx,
			log.NewEntry(log.StandardLogger()).
				WithFields(log.Fields{
					"chain.id": bc.ChainID().Text(10),
				}),
		)

		block, err := bc.ec.BlockByNumber(logCtx, bc.ChainID(), big.NewInt(blockNumber))
		if err != nil {
			bFuture.Err() <- &types.TxListenerError{
				ChainID: bc.ChainID(),
				Err:     errors.FromError(err).ExtendComponent(component),
			}
			return
		}

		// Block should be available
		if block == nil {
			bFuture.Err() <- &types.TxListenerError{
				ChainID: bc.ChainID(),
				Err:     errors.NotFoundError("block %v missing", blockNumber).ExtendComponent(component),
			}
			return
		}

		bl := &types.TxListenerBlock{
			ChainID:  bc.ChainID(),
			Block:    *block,
			Receipts: []*types.TxListenerReceipt{},
		}

		// Retrieve receipts in separate go-routines
		var rFutures []*types.Future
		for _, tx := range bl.Block.Transactions() {
			rFutures = append(rFutures, bc.fetchReceipt(ctx, tx.Hash()))
		}

		for i, rFuture := range rFutures {
			select {
			case err := <-rFuture.Err():
				bFuture.Err() <- &types.TxListenerError{
					ChainID: bc.ChainID(),
					Err:     errors.FromError(err).ExtendComponent(component),
				}
				return
			case res := <-rFuture.Result():
				receipt := res.(*types.TxListenerReceipt)
				receipt.TxIndex = uint64(i)
				receipt.BlockHash = block.Hash()
				receipt.BlockNumber = block.Number().Int64()
				bl.Receipts = append(bl.Receipts, receipt)
			}
		}

		// Return block in result
		bFuture.Result() <- bl
	}()

	return bFuture
}

// fetchReceipt fetch a receipt asynchronously
func (bc *BlockCursor) fetchReceipt(ctx context.Context, txHash common.Hash) *types.Future {
	future := types.NewFuture()

	log.WithFields(log.Fields{
		"tx.hash": txHash.Hex(),
	}).Tracef("tx-listener: fetch receipt")
	go func() {
		defer future.Close()
		logCtx := logger.WithLogEntry(
			ctx,
			log.NewEntry(log.StandardLogger()).
				WithFields(log.Fields{
					"chain.id": bc.ChainID().Text(10),
				}),
		)
		receipt, err := bc.ec.TransactionReceipt(logCtx, bc.ChainID(), txHash)
		if err != nil {
			future.Err() <- &types.TxListenerError{
				ChainID: bc.ChainID(),
				Err:     errors.FromError(err).ExtendComponent(component),
			}
			return
		}

		if receipt == nil {
			future.Err() <- &types.TxListenerError{
				ChainID: bc.ChainID(),
				Err:     errors.NotFoundError("receipt %q missing", txHash.Hex()).ExtendComponent(component),
			}
			return
		}

		r := &types.TxListenerReceipt{
			ChainID: bc.ChainID(),
			Receipt: *receipt,
			TxHash:  txHash,
		}
		future.Result() <- r
	}()

	return future
}
