package grpc

import (
	"context"

	"google.golang.org/grpc"

	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

type serviceBuilder struct {
	service svc.EnvelopeStoreServer
}

func newServiceBuilder(service svc.EnvelopeStoreServer) *serviceBuilder {
	return &serviceBuilder{
		service: service,
	}
}

func (b *serviceBuilder) Build(ctx context.Context, name string, configuration interface{}) (func(srv *grpc.Server), error) {
	registrator := func(srv *grpc.Server) {
		svc.RegisterEnvelopeStoreServer(srv, b.service)
	}

	return registrator, nil
}
