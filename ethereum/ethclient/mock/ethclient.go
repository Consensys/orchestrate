package mock

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/ethereum/types"
)

// Client is a Mock client
type Client struct {
	blocks map[string][]*ethtypes.Block

	mux  *sync.RWMutex
	head map[string]uint64
}

// NewClient creates a new mock client
func NewClient(blocks map[string][]*ethtypes.Block) *Client {
	head := make(map[string]uint64)
	for chain := range blocks {
		head[chain] = 0
	}
	return &Client{
		blocks: blocks,
		mux:    &sync.RWMutex{},
		head:   head,
	}
}

func (ec *Client) Mine(chainID *big.Int) {
	ec.mux.Lock()
	defer ec.mux.Unlock()

	if int(ec.head[chainID.Text(10)])+1 < len(ec.blocks[chainID.Text(10)]) {
		ec.head[chainID.Text(10)]++
	}
}

type ctxKeyType string

const errorCtxKey ctxKeyType = "error"

func WithError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errorCtxKey, err)
}

func GetError(ctx context.Context) error {
	err, ok := ctx.Value(errorCtxKey).(error)
	if !ok {
		return nil
	}
	return err
}

func (ec *Client) Networks(ctx context.Context) []*big.Int {
	chains := []*big.Int{}
	for chain := range ec.blocks {
		c, ok := big.NewInt(0).SetString(chain, 10)
		if !ok {
			panic("invalid chain id")
		}
		chains = append(chains, c)
	}
	return chains
}

func (ec *Client) BlockByNumber(ctx context.Context, chainID, number *big.Int) (*ethtypes.Block, error) {
	err := GetError(ctx)
	if err != nil {
		return nil, err
	}

	// Simulate io time
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(2 * time.Millisecond):
		ec.mux.RLock()
		defer ec.mux.RUnlock()

		if number == nil {
			number = big.NewInt(int64(ec.head[chainID.Text(10)]))
		}

		if number.Uint64() <= ec.head[chainID.Text(10)] {
			block := ec.blocks[chainID.Text(10)][number.Uint64()]
			header := ethtypes.CopyHeader(block.Header())
			header.Number = number
			blck := ethtypes.NewBlockWithHeader(header)
			return blck.WithBody(block.Transactions(), block.Uncles()), nil
		}

		if number.Uint64() > ec.head[chainID.Text(10)] {
			return nil, nil
		}
		return nil, fmt.Errorf("error")
	}
}

func (ec *Client) HeaderByNumber(ctx context.Context, chainID, number *big.Int) (*ethtypes.Header, error) {
	err := GetError(ctx)
	if err != nil {
		return nil, err
	}

	// Simulate io time
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(2 * time.Millisecond):
		ec.mux.RLock()
		defer ec.mux.RUnlock()
		if number == nil {
			number = big.NewInt(int64(ec.head[chainID.Text(10)]))
		}

		if number.Uint64() <= ec.head[chainID.Text(10)] {
			block := ec.blocks[chainID.Text(10)][number.Uint64()]
			header := ethtypes.CopyHeader(block.Header())
			header.Number = number
			return header, nil
		}

		if number.Uint64() > ec.head[chainID.Text(10)] {
			return nil, nil
		}

		return nil, fmt.Errorf("error")
	}
}

func (ec *Client) TransactionReceipt(ctx context.Context, chainID *big.Int, txHash ethcommon.Hash) (*ethtypes.Receipt, error) {
	err := GetError(ctx)
	if err != nil {
		return nil, err
	}

	// Simulate io time
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(2 * time.Millisecond):
		ec.mux.RLock()
		defer ec.mux.RUnlock()
		for _, block := range ec.blocks[chainID.Text(10)][:ec.head[chainID.Text(10)]+1] {
			if block.Transaction(txHash) != nil {
				return &ethtypes.Receipt{
					TxHash: txHash,
				}, nil
			}
		}
		return nil, nil
	}
}

func (ec *Client) CodeAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) ([]byte, error) {
	err := GetError(ctx)
	if err != nil {
		return nil, err
	}

	// Simulate io time
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(2 * time.Millisecond):
		ec.mux.RLock()
		defer ec.mux.RUnlock()

		if blockNumber == nil {
			blockNumber = big.NewInt(int64(ec.head[chainID.Text(10)]))
		}

		if blockNumber.Uint64() <= ec.head[chainID.Text(10)] {
			return []byte{1, 2, 3}, nil
		}

		if blockNumber.Uint64() > ec.head[chainID.Text(10)] {
			return nil, nil
		}
		return nil, fmt.Errorf("error")
	}
}

func (ec *Client) SyncProgress(ctx context.Context, chainID *big.Int) (*eth.SyncProgress, error) {
	return nil, fmt.Errorf("not implemented error")
}

func (ec *Client) TransactionByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (tx *ethtypes.Transaction, isPending bool, err error) {
	return nil, false, fmt.Errorf("not implemented error")
}

func (ec *Client) BlockByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (*ethtypes.Block, error) {
	return nil, fmt.Errorf("not implemented error")
}

func (ec *Client) HeaderByHash(ctx context.Context, chainID *big.Int, hash ethcommon.Hash) (*ethtypes.Header, error) {
	return nil, fmt.Errorf("not implemented error")
}

func (ec *Client) BalanceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) (*big.Int, error) {
	return nil, fmt.Errorf("not implemented error")
}

func (ec *Client) StorageAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, key ethcommon.Hash, blockNumber *big.Int) ([]byte, error) {
	return nil, fmt.Errorf("not implemented error")
}

func (ec *Client) NonceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, blockNumber *big.Int) (uint64, error) {
	return 0, fmt.Errorf("not implemented error")
}

func (ec *Client) PendingBalanceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) (*big.Int, error) {
	return nil, fmt.Errorf("not implemented error")
}

func (ec *Client) PendingStorageAt(ctx context.Context, chainID *big.Int, account ethcommon.Address, key ethcommon.Hash) ([]byte, error) {
	return nil, fmt.Errorf("not implemented error")
}

func (ec *Client) PendingCodeAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) ([]byte, error) {
	return nil, fmt.Errorf("not implemented error")
}

func (ec *Client) PendingNonceAt(ctx context.Context, chainID *big.Int, account ethcommon.Address) (uint64, error) {
	return 0, fmt.Errorf("not implemented error")
}

func (ec *Client) CallContract(ctx context.Context, chainID *big.Int, msg *eth.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return nil, fmt.Errorf("not implemented error")
}

func (ec *Client) PendingCallContract(ctx context.Context, chainID *big.Int, msg *eth.CallMsg) ([]byte, error) {
	return nil, fmt.Errorf("not implemented error")
}

func (ec *Client) SuggestGasPrice(ctx context.Context, chainID *big.Int) (*big.Int, error) {
	return nil, fmt.Errorf("not implemented error")
}

func (ec *Client) EstimateGas(ctx context.Context, chainID *big.Int, msg *eth.CallMsg) (uint64, error) {
	return 0, fmt.Errorf("not implemented error")
}

func (ec *Client) SendRawPrivateTransaction(ctx context.Context, chainID *big.Int, raw string, args *types.PrivateArgs) (ethcommon.Hash, error) {
	return [ethcommon.HashLength]byte{}, fmt.Errorf("not implemented error")
}

func (ec *Client) SendRawTransaction(ctx context.Context, chainID *big.Int, raw string) error {
	return fmt.Errorf("not implemented error")
}

func (ec *Client) SendTransaction(ctx context.Context, chainID *big.Int, args *types.SendTxArgs) (ethcommon.Hash, error) {
	return [ethcommon.HashLength]byte{}, fmt.Errorf("not implemented error")
}
