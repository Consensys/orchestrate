package http

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	grpcgateway "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler/grpc-gateway"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/proto"
)

func NewBuilder(
	service svc.ContractRegistryServer,
	opts ...runtime.ServeMuxOption,
) *grpcgateway.Builder {
	registrator := func(ctx context.Context, mux *runtime.ServeMux) error {
		return svc.RegisterContractRegistryHandlerServer(ctx, mux, service)
	}
	return grpcgateway.NewBuilder(
		opts,
		[]func(ctx context.Context, mux *runtime.ServeMux) error{
			registrator,
		},
	)

}
