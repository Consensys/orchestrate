package grpc

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	types "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/contract-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/contract-registry.git/app/grpc/services"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/contract-registry.git/app/infra"
)

// Variable use as GRPC server singleton for injection pattern
var (
	s         *server
	initOnce  = &sync.Once{}
	closeOnce = &sync.Once{}
	ready     = &atomic.Value{}
)

func init() {
	s = new()
}

type server struct {
	grpc *grpc.Server
}

// new create a new server
func new() *server {
	return &server{}
}

// Init inititalize grpc server
func Init() {
	initOnce.Do(func() {
		s.grpc = CreateServer()
		ready.Store(true)
		log.Infof("grpc: server ready")
	})
}

// Ready indicate if server is ready
func Ready() bool {
	return ready.Load().(bool)
}

// Server returns GRPC server
func Server() *grpc.Server {
	if !Ready() {
		log.Fatal("GRPC server is not ready. Please call Init() first")
	}
	return s.grpc
}

// Close GRPC server
func Close(ctx context.Context) {
	closeOnce.Do(func() {
		log.Debugf("grpc: closing...")
		Server().GracefulStop()
		log.Debugf("grpc: closed")
	})
}

// CreateServer creates grpc server
// CreateServer must called after infrastructure has been set
func CreateServer() *grpc.Server {
	// Set log entry
	logEntry := log.NewEntry(log.StandardLogger())
	grpc_logrus.ReplaceGrpcLogger(logEntry)

	panicHandler := grpc_recovery.RecoveryHandlerFunc(func(p interface{}) error {
		buf := make([]byte, 1<<16)
		runtime.Stack(buf, true)
		logEntry.Errorf("panic recovered: %+v", string(buf))
		return status.Errorf(codes.Internal, "%s", p)
	})

	// Set GRPC Interceptors
	server := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(infra.Tracer())),
			grpc_logrus.StreamServerInterceptor(logEntry),
			grpc_prometheus.StreamServerInterceptor,
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(panicHandler)),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(infra.Tracer())),
			grpc_logrus.UnaryServerInterceptor(logEntry),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(panicHandler)),
		)),
	)

	// Register services
	types.RegisterRegistryServer(server, services.NewRegistryService(infra.Registry()))

	// Register prometheus
	grpc_prometheus.Register(server)

	return server
}
