package ethclient

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// MultiEthClient is client that can connect to multiple Ethereum chains
type MultiEthClient struct {
	ecs map[string]*EthClient
}

// MultiDial connects a multi-client to the given URLs.
func MultiDial(rawurls []string) (*MultiEthClient, error) {
	return MultiDialContext(context.Background(), rawurls)
}

// MultiDialContext connects a multi-client to the given URLs.
func MultiDialContext(ctx context.Context, rawurls []string) (*MultiEthClient, error) {
	errors := make(chan error, len(rawurls))
	clients := make(chan *EthClient, len(rawurls))
	wait := &sync.WaitGroup{}
	for _, rawurl := range rawurls {
		wait.Add(1)
		go func(rawurl string) {
			defer wait.Done()
			c, err := DialContext(ctx, rawurl)
			if err != nil {
				errors <- err
			} else {
				clients <- c
			}
		}(rawurl)
	}
	wait.Wait()
	close(clients)
	close(errors)

	if len(errors) > 1 {
		return nil, <-errors
	}

	ecs := make(map[string]*EthClient)
	for ec := range clients {
		chainID, err := ec.NetworkID(ctx)
		if err != nil {
			panic(err)
		}
		ecs[chainIDToString(chainID)] = ec
	}

	return &MultiEthClient{ecs: ecs}, nil
}

// Networks return networks ID multi client is connected to
func (mec *MultiEthClient) Networks(ctx context.Context) []*big.Int {
	networks := []*big.Int{}
	for _, ec := range mec.ecs {
		chain, _ := ec.NetworkID(ctx)
		if chain != nil {
			networks = append(networks, chain)
		}
	}
	return networks
}

// HeaderByHash returns the block header with the given hash.
func (mec *MultiEthClient) HeaderByHash(ctx context.Context, chainID *big.Int, hash common.Hash) (*types.Header, error) {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return nil, fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.HeaderByHash(ctx, hash)
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (mec *MultiEthClient) HeaderByNumber(ctx context.Context, chainID *big.Int, number *big.Int) (*types.Header, error) {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return nil, fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.HeaderByNumber(ctx, number)
}

// BlockByNumber returns a block from the current canonical chain. If number is
// nil, the latest known header is returned.
func (mec *MultiEthClient) BlockByNumber(ctx context.Context, chainID *big.Int, number *big.Int) (*types.Block, error) {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return nil, fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.BlockByNumber(ctx, number)
}

// TransactionByHash returns the transaction with the given hash.
func (mec *MultiEthClient) TransactionByHash(ctx context.Context, chainID *big.Int, hash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return nil, false, fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.TransactionByHash(ctx, hash)
}

// TransactionCount returns the total number of transactions in the given block
func (mec *MultiEthClient) TransactionCount(ctx context.Context, chainID *big.Int, blockHash common.Hash) (uint, error) {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return 0, fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.TransactionCount(ctx, blockHash)
}

// TransactionReceipt returns the receipt of a transaction by transaction hash.
// Note that the receipt is not available for pending transactions.
func (mec *MultiEthClient) TransactionReceipt(ctx context.Context, chainID *big.Int, txHash common.Hash) (*types.Receipt, error) {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return nil, fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.TransactionReceipt(ctx, txHash)
}

// BalanceAt returns the wei balance of the given account.
// The block number can be nil, in which case the balance is taken from the latest known block.
func (mec *MultiEthClient) BalanceAt(ctx context.Context, chainID *big.Int, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return nil, fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.BalanceAt(ctx, account, blockNumber)
}

// PendingBalanceAt returns the wei balance of the given account in the pending state.
func (mec *MultiEthClient) PendingBalanceAt(ctx context.Context, chainID *big.Int, account common.Address) (*big.Int, error) {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return nil, fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.PendingBalanceAt(ctx, account)
}

// NonceAt returns the account nonce of the given account.
// The block number can be nil, in which case the nonce is taken from the latest known block.
func (mec *MultiEthClient) NonceAt(ctx context.Context, chainID *big.Int, account common.Address, blockNumber *big.Int) (uint64, error) {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return 0, fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.NonceAt(ctx, account, blockNumber)
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (mec *MultiEthClient) PendingNonceAt(ctx context.Context, chainID *big.Int, account common.Address) (uint64, error) {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return 0, fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.PendingNonceAt(ctx, account)
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (mec *MultiEthClient) SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error) {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return nil, fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.SuggestGasPrice(ctx)
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the current pending state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func (mec *MultiEthClient) EstimateGas(ctx context.Context, chainID *big.Int, msg ethereum.CallMsg) (uint64, error) {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return 0, fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.EstimateGas(ctx, msg)
}

// SendRawTransaction allows to send a raw transaction
func (mec *MultiEthClient) SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.SendRawTransaction(ctx, raw)
}

// SyncProgress retrieves client current progress of the sync algorithm.
func (mec *MultiEthClient) SyncProgress(ctx context.Context, chainID *big.Int) (*ethereum.SyncProgress, error) {
	ec, ok := mec.getClient(chainID)
	if !ok {
		return nil, fmt.Errorf("No client registered for %v", chainID)
	}
	return ec.SyncProgress(ctx)
}

func (mec *MultiEthClient) getClient(chainID *big.Int) (*EthClient, bool) {
	ec, ok := mec.ecs[chainIDToString(chainID)]
	return ec, ok
}

func chainIDToString(chainID *big.Int) string {
	return chainID.Text(16)
}
