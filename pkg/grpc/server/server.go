package grpcserver

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	grpcerror "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/grpc/error"
	"google.golang.org/grpc"
)

// NewServerWithDefaultOptions creates a new server with default options
func NewServerWithDefaultOptions() *grpc.Server {
	return grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
			grpc_logrus.StreamServerInterceptor(log.NewEntry(log.StandardLogger())),
			grpc_prometheus.StreamServerInterceptor,
			grpcerror.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(RecoverPanicHandler)),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
			grpc_logrus.UnaryServerInterceptor(log.NewEntry(log.StandardLogger())),
			grpc_prometheus.UnaryServerInterceptor,
			grpcerror.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(RecoverPanicHandler)),
		)),
	)
}
