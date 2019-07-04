package geth

import (
	"context"
	"fmt"
	"math/big"
)

type NullClient struct {
	chainID *big.Int
}

func CreateNullClient(chainID *big.Int) *NullClient {
	return &NullClient{
		chainID: chainID,
	}
}

func (c *NullClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return fmt.Errorf("no RPC connection registered for chain %q", c.chainID.String())
}

func (c *NullClient) Close() {
	// Do nothing
}
