package ethclient

import (
	"context"
	"math/big"

	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/types"
)

// TransactionSender is a service for sending transaction to a blockchain
type TransactionSender interface {
	// SendTransaction injects a signed transaction into the pending pool for execution.
	SendTransaction(ctx context.Context, chainID *big.Int, args *types.SendTxArgs) (ethcommon.Hash, error)

	// SendRawTransaction allows to send a raw transaction
	SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error

	// SendRawPrivateTransaction send a raw transaction to a Ethreum node supporting privacy (e.g Quorum+Tessera node)
	SendRawPrivateTransaction(ctx context.Context, chainID *big.Int, raw string, args *types.PrivateArgs) (ethcommon.Hash, error)
}

// ChainLedgerReader is a service to access a blockchain ledger information
type ChainLedgerReader interface {
	// BlockByHash returns the given full block.
	BlockByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (*ethtypes.Block, error)

	// BlockByNumber returns a block from the current canonical chain
	BlockByNumber(ctx context.Context, chainID *big.Int, number *big.Int) (*ethtypes.Block, error)

	// HeaderByHash returns the block header with the given hash.
	HeaderByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (*ethtypes.Header, error)

	// HeaderByNumber returns a block header from the current canonical chain. If number is
	// nil, the latest known header is returned.
	HeaderByNumber(ctx context.Context, chainID *big.Int, number *big.Int) (*ethtypes.Header, error)

	// TransactionByHash returns the transaction with the given hash.
	TransactionByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (tx *ethtypes.Transaction, isPending bool, err error)

	// TransactionReceipt returns the receipt of a transaction by transaction hash.
	TransactionReceipt(ctx context.Context, chainID *big.Int, txHash ethcommon.Hash) (*ethtypes.Receipt, error)
}

// ChainStateReader is a service to access a blockchain state information
type ChainStateReader interface {
	// BalanceAt returns wei balance of the given account.
	// The block number can be nil, in which case the balance is taken from the latest known block.
	BalanceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) (*big.Int, error)

	// StorageAt returns value of key in the contract storage of the given account.
	// The block number can be nil, in which case the value is taken from the latest known block.
	StorageAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, key ethcommon.Hash, blockNumber *big.Int) ([]byte, error)

	// CodeAt returns contract code of the given account.
	// The block number can be nil, in which case the code is taken from the latest known block.
	CodeAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) ([]byte, error)

	// NonceAt returns account nonce of the given account.
	// The block number can be nil, in which case the nonce is taken from the latest known block.
	NonceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) (uint64, error)

	// PendingBalanceAt returns wei balance of the given account in the pending state.
	PendingBalanceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) (*big.Int, error)

	// PendingStorageAt returns value of key in the contract storage of the given account in the pending state.
	PendingStorageAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, key ethcommon.Hash) ([]byte, error)

	// PendingCodeAt returns contract code of the given account in the pending state.
	PendingCodeAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) ([]byte, error)

	// PendingNonceAt returns account nonce of the given account in the pending state.
	PendingNonceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) (uint64, error)
}

// ContractCaller is a service to perform contract calls
type ContractCaller interface {
	// CallContract executes a message call transaction, which is directly executed in the VM
	CallContract(ctx context.Context, chainID *big.Int, msg *eth.CallMsg, blockNumber *big.Int) ([]byte, error)

	// PendingCallContract executes a message call transaction using the EVM.
	// The state seen by the contract call is the pending state.
	PendingCallContract(ctx context.Context, chainID *big.Int, msg *eth.CallMsg) ([]byte, error)
}

// GasEstimator is a service that can provide transaction gas price estimation
type GasEstimator interface {
	// EstimateGas tries to estimate the gas needed to execute a specific transaction
	EstimateGas(ctx context.Context, chainID *big.Int, msg *eth.CallMsg) (uint64, error)
}

// GasPricer is service that
type GasPricer interface {
	// SuggestGasPrice retrieves the currently suggested gas price
	SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error)
}

// ChainSyncReader is a service to access to the node's current sync status
type ChainSyncReader interface {
	SyncProgress(ctx context.Context, chainID *big.Int) (*eth.SyncProgress, error)
}

type Client interface {
	TransactionSender
	ChainLedgerReader
	ChainStateReader
	ContractCaller
	GasEstimator
	GasPricer
	ChainSyncReader
}
