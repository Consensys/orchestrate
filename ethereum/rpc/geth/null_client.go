package geth

import (
	"context"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
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
	return errors.EthConnectionError("no RPC connection registered for chain %q", c.chainID.Text(10))
}

func (c *NullClient) Close() {
	// Do nothing
}
