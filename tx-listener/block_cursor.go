package listener

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

// TxListenerError is what is provided to the user when an error occurs.
type TxListenerError struct {
	ChainID *big.Int
	Err     error
}

// Progress stores information about current cursor position listening to a chain
type Progress struct {
	CurrentBlock uint64 // Current block number where sync is at
	HighestBlock uint64 // Highest available block
}

// BlockCursor allows to retrieve a new block
type BlockCursor interface {
	ChainID() *big.Int
	Next(ctx context.Context) (*types.Block, error)
	Set(blockNumber uint64)
	Progress(ctx context.Context) *Progress
	Close()
}

// TxListenerEthClient is a minimal EthClient interface required by a TxListener
type TxListenerEthClient interface {
	// BlockByNumber retrieve a block by its number
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)

	// SyncProgress retrieves client current progress of the sync algorithm.
	SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error)

	// NetworkID returns the network ID
	NetworkID(ctx context.Context) (*big.Int, error)
}

type blockCursor struct {
	ec TxListenerEthClient

	chainID *big.Int

	mux     *sync.RWMutex
	pos     uint64
	current *types.Block

	// TODO: add an history of blocks of configurable length
	// (we could add some checks to ensure no re-org happened)
}

func (bc *blockCursor) ChainID() *big.Int {
	return bc.chainID
}

func (bc *blockCursor) Set(blockNumber uint64) {
	bc.mux.Lock()
	defer bc.mux.Unlock()
	bc.pos = blockNumber
}

func (bc *blockCursor) Next(ctx context.Context) (*types.Block, error) {
	block, err := bc.ec.BlockByNumber(ctx, big.NewInt(int64(bc.pos+1)))
	if err != nil {
		return nil, err
	}

	if block != nil {
		bc.mux.Lock()
		if block.NumberU64() == bc.pos+1 {
			// Position has not been changed
			bc.pos = bc.pos + 1
		}
		bc.mux.Unlock()
	}

	return block, nil
}

func newBlockCursor(ec TxListenerEthClient, chainID *big.Int) (*blockCursor, error) {
	cursor := &blockCursor{
		ec:      ec,
		mux:     &sync.Mutex{},
		chainID: chainID,
	}
	return cursor, nil
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
