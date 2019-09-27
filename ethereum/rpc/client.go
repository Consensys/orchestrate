package rpc

import (
	"context"
)

// Client interface for an Geth RPC Client
type Client interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
	Close()
}
