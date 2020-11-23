package tags

import (
	"context"
	"fmt"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/config/static"
	"google.golang.org/grpc"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, func(srv *grpc.Server), error) {
	cfg, ok := configuration.(*static.Tags)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid interceptor configuration type (expected %T but got %T)", cfg, configuration)
	}

	opt := grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)

	return grpc_ctxtags.UnaryServerInterceptor(opt), grpc_ctxtags.StreamServerInterceptor(opt), nil, nil
}
