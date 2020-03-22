package grpcrecovery

import (
	"context"
	"fmt"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"google.golang.org/grpc"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(ctx context.Context, name string, configuration interface{}) (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, func(srv *grpc.Server), error) {
	cfg, ok := configuration.(*static.Recovery)
	if !ok {
		return nil, nil, nil, fmt.Errorf("invalid interceptor configuration type (expected %T but got %T)", cfg, configuration)
	}
	return UnaryServerInterceptor(), StreamServerInterceptor(), nil, nil
}

// UnaryServerInterceptor returns a grpc unary server interceptor (middleware) that allows
// to intercept internal errors
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(RecoverPanicHandler))
}

// StreamServerInterceptor returns a grpc streaming server interceptor for panic recovery.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(RecoverPanicHandler))
}
