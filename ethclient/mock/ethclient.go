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
