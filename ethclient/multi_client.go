package ethclient

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	eth "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// MultiClient is a client that can connect to multiple Ethereum chains
type MultiClient struct {
	mux *sync.Mutex
	ecs map[string]*EthClient
}

// NewMultiClient creates a new MultiClient
func NewMultiClient() *MultiClient {
	return &MultiClient{
		mux: &sync.Mutex{},
		ecs: make(map[string]*EthClient),
	}
}

// Dial an Ethereum client
func (mec *MultiClient) Dial(ctx context.Context, rawurl string) error {
	// Dial Ethereum client
	ec, err := DialContext(ctx, rawurl)
	if err != nil {
		return err
	}

	// Retrieve NetworkID
	chainID, err := ec.NetworkID(ctx)
	if err != nil {
		return err
	}

	// Register client
	mec.mux.Lock()
	mec.ecs[ChainIDToString(chainID)] = ec
	mec.mux.Unlock()

	return nil
}

// MultiDial connects to multiple Ethereum clients concurrently
func (mec *MultiClient) MultiDial(ctx context.Context, rawurls []string) error {
	// Dial clients concurrently
	wait := &sync.WaitGroup{}
	errors := make(chan error, len(rawurls))
	for _, rawurl := range rawurls {
		wait.Add(1)
		go func(rawurl string) {
			err := mec.Dial(ctx, rawurl)
			if err != nil {
				errors <- err
			}
			wait.Done()
		}(rawurl)
	}

	// Wait for all clients to be ready
	wait.Wait()
	close(errors)

	// In case we failed to connect to a client we return an error
	for err := range errors {
		return err
	}

	return nil
}

// Networks return networks ID multi client is connected to
func (mec *MultiClient) Networks(ctx context.Context) []*big.Int {
	networks := []*big.Int{}
	for _, ec := range mec.ecs {
		chain, _ := ec.NetworkID(ctx)
		if chain != nil {
			networks = append(networks, chain)
		}
	}
	return networks
}

// Close multiclient
func (mec *MultiClient) Close() {
	for _, ec := range mec.ecs {
		ec.Close()
	}
}

func (mec *MultiClient) getClient(chainID *big.Int) (*EthClient, error) {
	ec, ok := mec.ecs[ChainIDToString(chainID)]
	if !ok {
		return nil, fmt.Errorf("no client registered for chain %q", ChainIDToString(chainID))
	}
	return ec, nil
}

// HeaderByHash returns the block header with the given hash.
func (mec *MultiClient) HeaderByHash(ctx context.Context, chainID *big.Int, hash common.Hash) (*types.Header, error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return nil, err
	}
	return ec.HeaderByHash(ctx, hash)
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (mec *MultiClient) HeaderByNumber(ctx context.Context, chainID, number *big.Int) (*types.Header, error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return nil, err
	}
	return ec.HeaderByNumber(ctx, number)
}

// BlockByNumber returns a block from the current canonical chain. If number is
// nil, the latest known header is returned.
func (mec *MultiClient) BlockByNumber(ctx context.Context, chainID, number *big.Int) (*types.Block, error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return nil, err
	}
	return ec.BlockByNumber(ctx, number)
}

// TransactionByHash returns the transaction with the given hash.
func (mec *MultiClient) TransactionByHash(ctx context.Context, chainID *big.Int, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return nil, false, err
	}
	return ec.TransactionByHash(ctx, hash)
}

// TransactionCount returns the total number of transactions in the given block
func (mec *MultiClient) TransactionCount(ctx context.Context, chainID *big.Int, blockHash common.Hash) (uint, error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return 0, err
	}
	return ec.TransactionCount(ctx, blockHash)
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (mec *MultiClient) TransactionReceipt(ctx context.Context, chainID *big.Int, txHash common.Hash) (*types.Receipt, error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return nil, err
	}
	return ec.TransactionReceipt(ctx, txHash)
}

// BalanceAt returns the wei balance of the given account.
// The block number can be nil, in which case the balance is taken from the latest known block.
func (mec *MultiClient) BalanceAt(ctx context.Context, chainID *big.Int, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return nil, err
	}
	return ec.BalanceAt(ctx, account, blockNumber)
}

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (mec *MultiClient) PendingBalanceAt(ctx context.Context, chainID *big.Int, account common.Address) (*big.Int, error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return nil, err
	}
	return ec.PendingBalanceAt(ctx, account)
}

// NonceAt returns the account nonce of the given account.
// The block number can be nil, in which case the nonce is taken from the latest known block.
func (mec *MultiClient) NonceAt(ctx context.Context, chainID *big.Int, account common.Address, blockNumber *big.Int) (uint64, error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return 0, err
	}
	return ec.NonceAt(ctx, account, blockNumber)
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (mec *MultiClient) PendingNonceAt(ctx context.Context, chainID *big.Int, account common.Address) (uint64, error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return 0, err
	}
	return ec.PendingNonceAt(ctx, account)
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (mec *MultiClient) SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return nil, err
	}
	return ec.SuggestGasPrice(ctx)
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the current pending state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func (mec *MultiClient) EstimateGas(ctx context.Context, chainID *big.Int, msg *eth.CallMsg) (uint64, error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return 0, err
	}
	return ec.EstimateGas(ctx, *msg)
}

// SyncProgress retrieves client current progress of the sync algorithm.
func (mec *MultiClient) SyncProgress(ctx context.Context, chainID *big.Int) (*eth.SyncProgress, error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return nil, err
	}
	return ec.SyncProgress(ctx)
}

// TxSender is an interface for a transaction sender compatible with Ethereum nodes supporting privacy
type TxSender interface {
	// SendRawTransaction allows to send a raw transaction to an Ethereum node
	SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error

	// SendTransaction send a transaction to Ethereum node
	// Should be compatible with Ethereum nodes supporting privacy (such as Quorum)
	SendTransaction(ctx context.Context, chainID *big.Int, args *SendTxArgs) (txHash common.Hash, err error)

	// SendRawPrivateTransaction send a raw transaction to a Ethereum node supporting privacy (e.g Quorum+Tessera node)
	SendRawPrivateTransaction(ctx context.Context, chainID *big.Int, raw string, args *PrivateArgs) (common.Hash, error)
}

// SendRawTransaction allows to send a raw transaction
func (mec *MultiClient) SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return err
	}
	return ec.SendRawTransaction(ctx, raw)
}

// SendTransaction send a transaction to Ethereum node
func (mec *MultiClient) SendTransaction(ctx context.Context, chainID *big.Int, args *SendTxArgs) (txHash common.Hash, err error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return common.Hash{}, err
	}
	return ec.SendTransaction(ctx, args)
}

// SendRawPrivateTransaction send a raw transaction to a Ethereum node supporting privacy (e.g Quorum+Tessera node)
// TODO: to be implemented
func (mec *MultiClient) SendRawPrivateTransaction(ctx context.Context, chainID *big.Int, raw string, args *PrivateArgs) (common.Hash, error) {
	ec, err := mec.getClient(chainID)
	if err != nil {
		return common.Hash{}, err
	}
	return ec.SendRawPrivateTransaction(ctx, raw, args)
}
