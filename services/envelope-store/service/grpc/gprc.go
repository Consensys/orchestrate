package gprc

import (
	"context"

	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
	"google.golang.org/grpc"
)

type Builder struct {
	service svc.EnvelopeStoreServer
}

func NewBuilder(service svc.EnvelopeStoreServer) *Builder {
	return &Builder{
		service: service,
	}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (func(srv *grpc.Server), error) {
	registrator := func(srv *grpc.Server) {
		svc.RegisterEnvelopeStoreServer(srv, b.service)
	}

	return registrator, nil
}
