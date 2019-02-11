package listener

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EthClient is a minimal Ethereum Client interface required by a TxListener
type EthClient interface {
	// BlockByNumber retrieve a block by its number
	BlockByNumber(ctx context.Context, chainID *big.Int, number *big.Int) (*types.Block, error)

	// TransactionReceipt retrieve a transaction receipt using its hash
	TransactionReceipt(ctx context.Context, chainID *big.Int, txHash common.Hash) (*types.Receipt, error)

	// SyncProgress retrieves client current progress of the sync algorithm.
	SyncProgress(ctx context.Context, chainID *big.Int) (*ethereum.SyncProgress, error)
}
