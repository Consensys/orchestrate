package infra

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// EthClient embed a go-ethereum ethclient so we can implement SendRawTransaction
type EthClient struct {
	*ethclient.Client
	rpc *rpc.Client
}

// NewClient creates a client that uses the given RPC client.
func NewClient(c *rpc.Client) *EthClient {
	ec := ethclient.NewClient(c)
	return &EthClient{ec, c}
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (*EthClient, error) {
	return DialContext(context.Background(), rawurl)
}

// DialContext connects a client to the given URL.
func DialContext(ctx context.Context, rawurl string) (*EthClient, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

// SendRawTransaction allows to send a raw transaction
func (ec *EthClient) SendRawTransaction(ctx context.Context, raw string) error {
	return ec.rpc.CallContext(ctx, nil, "eth_sendRawTransaction", raw)
}
