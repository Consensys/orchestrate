package grpcauth

import (
	"context"
	"fmt"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"google.golang.org/grpc"
)

type Builder struct {
	checker      auth.Checker
	multitenancy bool
}

func NewBuilder(checker auth.Checker, multitenancy bool) *Builder {
	return &Builder{
		checker:      checker,
		multitenancy: multitenancy,
	}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, func(srv *grpc.Server), error) {
	cfg, ok := configuration.(*static.Auth)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid interceptor configuration type (expected %T but got %T)", cfg, configuration)
	}
	return UnaryServerInterceptor(b.checker, b.multitenancy), StreamServerInterceptor(b.checker, b.multitenancy), nil, nil
}

// UnaryServerInterceptor returns a grpc unary server interceptor (middleware) that allows
// to intercept internal errors
func UnaryServerInterceptor(checker auth.Checker, multitenancy bool) grpc.UnaryServerInterceptor {
	return grpc_auth.UnaryServerInterceptor(Auth(checker, multitenancy))
}

// StreamServerInterceptor returns a grpc streaming server interceptor for panic recovery.
func StreamServerInterceptor(checker auth.Checker, multitenancy bool) grpc.StreamServerInterceptor {
	return grpc_auth.StreamServerInterceptor(Auth(checker, multitenancy))
}
