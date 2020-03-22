package mock

import (
	"context"

	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store/memory"
	"google.golang.org/grpc"
)

// EnvelopeStoreClient is a client that wraps an EnvelopeStoreServer into an EnvelopeStoreClient
type EnvelopeStoreClient struct {
	srv svc.EnvelopeStoreServer
}

func New() *EnvelopeStoreClient {
	return &EnvelopeStoreClient{
		srv: memory.New(),
	}
}

func (client *EnvelopeStoreClient) Store(ctx context.Context, in *svc.StoreRequest, opts ...grpc.CallOption) (*svc.StoreResponse, error) {
	return client.srv.Store(ctx, in)
}

// Load envelope by identifier
func (client *EnvelopeStoreClient) LoadByID(ctx context.Context, in *svc.LoadByIDRequest, opts ...grpc.CallOption) (*svc.StoreResponse, error) {
	return client.srv.LoadByID(ctx, in)
}

// Load Envelope by transaction hash
func (client *EnvelopeStoreClient) LoadByTxHash(ctx context.Context, in *svc.LoadByTxHashRequest, opts ...grpc.CallOption) (*svc.StoreResponse, error) {
	return client.srv.LoadByTxHash(ctx, in)
}

// SetStatus set an envelope status
func (client *EnvelopeStoreClient) SetStatus(ctx context.Context, in *svc.SetStatusRequest, opts ...grpc.CallOption) (*svc.StatusResponse, error) {
	return client.srv.SetStatus(ctx, in)

}

// LoadPending load envelopes of pending transactions
func (client *EnvelopeStoreClient) LoadPending(ctx context.Context, in *svc.LoadPendingRequest, opts ...grpc.CallOption) (*svc.LoadPendingResponse, error) {
	return client.srv.LoadPending(ctx, in)

}
