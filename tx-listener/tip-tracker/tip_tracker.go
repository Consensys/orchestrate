package tiptracker

import (
	"context"
	"math/big"
)

// TipTracker keep track of blockchain highest final block
type TipTracker interface {
	// Return ID of the chain being tracked
	ChainID() *big.Int

	// Return block number of the highest final block in the canonical chain
	HighestBlock(ctx context.Context) (int64, error)
}
