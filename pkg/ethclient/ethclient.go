package ethclient

import (
	"context"
	"encoding/json"
	"math/big"

	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethereum/types"
	proto "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/ethereum"
)

//go:generate mockgen -source=ethclient.go -destination=mock/mock.go -package=mock

// TransactionSender is a service for sending transaction to a blockchain
type TransactionSender interface {
	// SendTransaction injects a signed transaction into the pending pool for execution.
	SendTransaction(ctx context.Context, url string, args *types.SendTxArgs) (ethcommon.Hash, error)

	// SendRawTransaction allows to send a raw transaction
	SendRawTransaction(ctx context.Context, url string, raw string) (ethcommon.Hash, error)

	// SendRawPrivateTransaction send a raw transaction to a Ethereum node supporting privacy with EEA privacy extensions
	SendRawPrivateTransaction(ctx context.Context, url string, raw string) (ethcommon.Hash, error)
}

type EEATransactionSender interface {
	// PrivDistributeRawTransaction Returns the enclaveKey of sent private transaction
	PrivDistributeRawTransaction(ctx context.Context, endpoint, raw string) (ethcommon.Hash, error)
}

type QuorumTransactionSender interface {
	// SendQuorumRawPrivateTransaction sends a raw signed transaction to a Quorum node
	// signedTxHash - is a hash returned by Quorum and then signed by a client
	// privateFor - is a list of public keys of Quorum nodes that can receive a private transaction
	SendQuorumRawPrivateTransaction(ctx context.Context, url string, signedTxHash string, privateFor []string) (ethcommon.Hash, error)

	// StoreRaw stores "data" field of a transaction in Tessera privacy enclave
	// It returns a hash of a stored transaction that should be used instead of transaction data
	StoreRaw(ctx context.Context, endpoint string, data []byte, privateFrom string) (string, error)
}

// ChainLedgerReader is a service to access a blockchain ledger information
type ChainLedgerReader interface {
	// BlockByHash returns the given full block.
	BlockByHash(ctx context.Context, url string, hash ethcommon.Hash) (*ethtypes.Block, error)

	// BlockByNumber returns a block from the current canonical chain
	BlockByNumber(ctx context.Context, url string, number *big.Int) (*ethtypes.Block, error)

	// HeaderByHash returns the block header with the given hash.
	HeaderByHash(ctx context.Context, url string, hash ethcommon.Hash) (*ethtypes.Header, error)

	// HeaderByNumber returns a block header from the current canonical chain. If number is
	// nil, the latest known header is returned.
	HeaderByNumber(ctx context.Context, url string, number *big.Int) (*ethtypes.Header, error)

	// TransactionByHash returns the transaction with the given hash.
	TransactionByHash(ctx context.Context, url string, hash ethcommon.Hash) (tx *ethtypes.Transaction, isPending bool, err error)

	// TransactionReceipt returns the receipt of a transaction by transaction hash.
	TransactionReceipt(ctx context.Context, url string, txHash ethcommon.Hash) (*proto.Receipt, error)
}

type EEAChainLedgerReader interface {
	// TransactionReceipt returns the receipt of a transaction by transaction hash.
	PrivateTransactionReceipt(ctx context.Context, url string, txHash ethcommon.Hash) (*proto.Receipt, error)
}

// ChainStateReader is a service to access a blockchain state information
type ChainStateReader interface {
	// BalanceAt returns wei balance of the given account.
	// The block number can be nil, in which case the balance is taken from the latest known block.
	BalanceAt(ctx context.Context, url string, account ethcommon.Address, blockNumber *big.Int) (*big.Int, error)

	// StorageAt returns value of key in the contract storage of the given account.
	// The block number can be nil, in which case the value is taken from the latest known block.
	StorageAt(ctx context.Context, url string, account ethcommon.Address, key ethcommon.Hash, blockNumber *big.Int) ([]byte, error)

	// CodeAt returns contract code of the given account.
	// The block number can be nil, in which case the code is taken from the latest known block.
	CodeAt(ctx context.Context, url string, account ethcommon.Address, blockNumber *big.Int) ([]byte, error)

	// NonceAt returns account nonce of the given account.
	// The block number can be nil, in which case the nonce is taken from the latest known block.
	NonceAt(ctx context.Context, url string, account ethcommon.Address, blockNumber *big.Int) (uint64, error)

	// PendingBalanceAt returns wei balance of the given account in the pending state.
	PendingBalanceAt(ctx context.Context, url string, account ethcommon.Address) (*big.Int, error)

	// PendingStorageAt returns value of key in the contract storage of the given account in the pending state.
	PendingStorageAt(ctx context.Context, url string, account ethcommon.Address, key ethcommon.Hash) ([]byte, error)

	// PendingCodeAt returns contract code of the given account in the pending state.
	PendingCodeAt(ctx context.Context, url string, account ethcommon.Address) ([]byte, error)

	// PendingNonceAt returns account nonce of the given account in the pending state.
	PendingNonceAt(ctx context.Context, url string, account ethcommon.Address) (uint64, error)
}

type EEAChainStateReader interface {
	// PrivEEANonce Returns the private transaction count for specified account and privacy group
	PrivEEANonce(ctx context.Context, endpoint string, account ethcommon.Address, privateFrom string, privateFor []string) (uint64, error)

	// PrivNonce Returns the private transaction count for specified account and privacy group
	PrivNonce(ctx context.Context, endpoint string, account ethcommon.Address, privacyGroupID string) (uint64, error)

	// EEAPrivPrecompiledContractAddr Returns the private precompiled contract address of Besu/Orion
	EEAPrivPrecompiledContractAddr(ctx context.Context, endpoint string) (ethcommon.Address, error)
}

// ChainStateReader is a service to access a blockchain state information
type QuorumChainStateReader interface {
	// GetStatus returns status of Tessera enclave if it is up or an error if it is down
	GetStatus(ctx context.Context, endpoint string) (status string, err error)
}

// ContractCaller is a service to perform contract calls
type ContractCaller interface {
	// CallContract executes a message call transaction, which is directly executed in the VM
	CallContract(ctx context.Context, url string, msg *eth.CallMsg, blockNumber *big.Int) ([]byte, error)

	// PendingCallContract executes a message call transaction using the EVM.
	// The state seen by the contract call is the pending state.
	PendingCallContract(ctx context.Context, url string, msg *eth.CallMsg) ([]byte, error)
}

// GasEstimator is a service that can provide transaction gas price estimation
type GasEstimator interface {
	// EstimateGas tries to estimate the gas needed to execute a specific transaction
	EstimateGas(ctx context.Context, url string, msg *eth.CallMsg) (uint64, error)
}

// GasPricer is service that
type GasPricer interface {
	// SuggestGasPrice retrieves the currently suggested gas price
	SuggestGasPrice(ctx context.Context, url string) (*big.Int, error)
}

// ChainSyncReader is a service to access to the node's current sync status
type ChainSyncReader interface {
	Network(ctx context.Context, url string) (*big.Int, error)
	SyncProgress(ctx context.Context, url string) (*eth.SyncProgress, error)
}

type Client interface {
	TransactionSender
	ChainLedgerReader
	ChainStateReader
	ContractCaller
	GasEstimator
	GasPricer
	ChainSyncReader
	Call(ctx context.Context, endpoint string, processResult func(result json.RawMessage) error, method string, args ...interface{}) error
}

type EEAClient interface {
	EEATransactionSender
	EEAChainLedgerReader
	EEAChainStateReader
	ContractCaller
	GasEstimator
	GasPricer
	ChainSyncReader
}

type QuorumClient interface {
	QuorumTransactionSender
	ChainLedgerReader
	QuorumChainStateReader
	ContractCaller
	GasEstimator
	GasPricer
	ChainSyncReader
}
