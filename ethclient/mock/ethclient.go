package mock

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Client is a Mock client
type Client struct {
	blocks []*ethtypes.Block

	mux  *sync.RWMutex
	head uint64
}

func NewClient(blocks []*ethtypes.Block) *Client {
	return &Client{
		blocks: blocks,
		mux:    &sync.RWMutex{},
	}
}

func (ec *Client) Mine() {
	ec.mux.Lock()
	defer ec.mux.Unlock()

	if int(ec.head)+1 < len(ec.blocks) {
		ec.head++
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
			number = big.NewInt(int64(ec.head))
		}

		if number.Uint64() <= ec.head {
			block := ec.blocks[number.Uint64()]
			header := ethtypes.CopyHeader(block.Header())
			header.Number = number
			blck := ethtypes.NewBlockWithHeader(header)
			return blck.WithBody(block.Transactions(), block.Uncles()), nil
		}

		if number.Uint64() > ec.head {
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
			number = big.NewInt(int64(ec.head))
		}

		if number.Uint64() <= ec.head {
			block := ec.blocks[number.Uint64()]
			header := ethtypes.CopyHeader(block.Header())
			header.Number = number
			return header, nil
		}

		if number.Uint64() > ec.head {
			return nil, nil
		}

		return nil, fmt.Errorf("error")
	}
}

func (ec *Client) TransactionReceipt(ctx context.Context, chainID *big.Int, txHash common.Hash) (*ethtypes.Receipt, error) {
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
		for _, block := range ec.blocks[:ec.head+1] {
			if block.Transaction(txHash) != nil {
				return &ethtypes.Receipt{
					TxHash: txHash,
				}, nil
			}
		}
		return nil, nil
	}
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
