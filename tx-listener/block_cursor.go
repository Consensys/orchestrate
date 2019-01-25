package listener

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

// BlockCursor allows to retrieve a new block
type BlockCursor interface {
	Next() (*types.Block, error)
	Set(pos *big.Int)
}

type blockCursor struct {
	c TxListenerEthClient

	next *big.Int
	// TODO: add an history of blocks of configurable length
	// (we could add some checks to ensure no re-org happened)
}

func newBlockCursor(c TxListenerEthClient) *blockCursor {
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
