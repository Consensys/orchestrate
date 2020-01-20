package grpcserver

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	grpcerror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/error"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func CodeToLevel(code codes.Code) log.Level {
	switch code {
	case codes.OK:
		return log.DebugLevel
	case codes.NotFound:
		return log.DebugLevel
	default:
		return grpc_logrus.DefaultCodeToLevel(code)
	}
}

// NewServer creates a new server with specific logrus options
func NewServer(auth authentication.Auth) *grpc.Server {

	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(CodeToLevel),
	}

	authF := Auth(
		auth,
		viper.GetBool(multitenancy.EnabledViperKey),
	)

	return grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
			grpc_logrus.StreamServerInterceptor(log.NewEntry(log.StandardLogger()), opts...),
			grpc_prometheus.StreamServerInterceptor,
			grpc_auth.StreamServerInterceptor(authF),
			grpcerror.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(RecoverPanicHandler)),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(opentracing.GlobalTracer())),
			grpc_logrus.UnaryServerInterceptor(log.NewEntry(log.StandardLogger()), opts...),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_auth.UnaryServerInterceptor(authF),
			grpcerror.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(RecoverPanicHandler)),
		)),
	)
}
