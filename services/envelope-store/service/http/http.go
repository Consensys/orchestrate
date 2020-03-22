package http

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	grpcgateway "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/grpc-gateway"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

func NewBuilder(
	service svc.EnvelopeStoreServer,
	opts ...runtime.ServeMuxOption,
) *grpcgateway.Builder {
	registrator := func(ctx context.Context, mux *runtime.ServeMux) error {
		return svc.RegisterEnvelopeStoreHandlerServer(ctx, mux, service)
	}
	return grpcgateway.NewBuilder(
		opts,
		[]func(ctx context.Context, mux *runtime.ServeMux) error{
			registrator,
		},
	)

}
