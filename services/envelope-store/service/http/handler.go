package http

import (
	"context"
	"reflect"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler"
	dynhandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/dynamic"
	grpcgateway "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/grpc-gateway"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

type handlerBuilder handler.Builder

func newHandlerBuilder(srv svc.EnvelopeStoreServer) handlerBuilder {
	builder := dynhandler.NewBuilder()

	// Add Builder for Contract-Registry API
	builder.AddBuilder(
		reflect.TypeOf(&dynamic.Envelopes{}),
		newGatewayBuilder(srv),
	)

	return builder
}

func newGatewayBuilder(
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
