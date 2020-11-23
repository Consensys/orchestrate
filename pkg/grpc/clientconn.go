package grpc

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/credentials"
	grpcerror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/grpc/interceptor/error"
	"google.golang.org/grpc"
)

// DialContext creates a client connection to the given target
func DialContext(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	conn, err = grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, errors.GRPCConnectionError(err.Error())
	}
	return conn, nil
}

// DialContextWithDefaultOptions creates a client with a set of default options
func DialContextWithDefaultOptions(ctx context.Context, target string) (conn *grpc.ClientConn, err error) {
	return DialContext(
		ctx,
		target,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&credentials.PerRPCCredentials{}),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpcerror.UnaryClientInterceptor(),
			grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
			grpc_logrus.UnaryClientInterceptor(log.NewEntry(log.StandardLogger())),
		)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			grpcerror.StreamClientInterceptor(),
			grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
			grpc_logrus.StreamClientInterceptor(log.NewEntry(log.StandardLogger())),
		)),
	)
}
