package client

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

func DialContextWithDefaultOptions(ctx context.Context, url string) (svc.EnvelopeStoreClient, error) {
	conn, err := grpc.DialContextWithDefaultOptions(
		ctx,
		url,
	)
	if err != nil {
		return nil, err
	}

	return svc.NewEnvelopeStoreClient(conn), nil
}
