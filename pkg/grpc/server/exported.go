package grpcserver

import (
	"context"
	"sync"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/errors"
	grpclogger "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/grpc/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/pkg/tracing/opentracing/jaeger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const component = "grpc.server"

var (
	initOnce   = &sync.Once{}
	server     *grpc.Server
	cmuxserver *CMuxServer
	enhancers  = []Enhancer{}
)

// Init initialize global GRPC server
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if server == nil {
			// Initialize opentracing tracer
			jaeger.Init(ctx)

			// Declare server with interceptors
			server = NewServerWithDefaultOptions()

			// Apply enhancers on server
			ApplyEnhancers(server, enhancers...)

			// Register server for prometheus metrics
			grpc_prometheus.Register(server)

			// Replace internal GRPC logger with a logrus logger
			grpclog.SetLoggerV2(
				&grpclogger.LogEntry{
					Entry: log.WithFields(log.Fields{"system": "grpc.internal"}),
				},
			)
		}

		// Log registered services
		services := []string{}
		for name := range server.GetServiceInfo() {
			services = append(services, name)
		}

		log.WithFields(log.Fields{
			"grpc.services": services,
		}).Infof("grpc: server ready")
	})
}

// GlobalServer return global GRPC server
func GlobalServer() *grpc.Server {
	return server
}

// SetGlobalServer sets global GRPC server
func SetGlobalServer(s *grpc.Server) {
	server = s
}

// AddEnhancers adds GRPC server enhancers that will be called at Init time
// Note that it should be called before Init()
func AddEnhancers(fns ...Enhancer) {
	enhancers = append(enhancers, fns...)
}

// ListenAndServe starts global server
func ListenAndServe() error {
	// Ensure GRPC server has been initialized
	if server == nil {
		log.Fatalf("grpc.server: GRPC server is not initialized")
	}

	// Ensure HTTP server has been initialized
	if http.GlobalServer() == nil {
		log.Fatalf("grpc.server: HTTP server is not initialized")
	}

	// Declare multiplexer server
	cmuxserver = NewCMuxServer(server, http.GlobalServer())

	// Serve requests
	err := cmuxserver.ListenAndServe("tcp", viper.GetString("http.hostname"))
	if err != nil {
		log.WithError(errors.FromError(err).ExtendComponent(component)).
			WithFields(log.Fields{
				"http.hostname": viper.GetString("http.hostname"),
			}).
			Error("grpc.server: error listening tcp connections")
		return err
	}

	log.Info("grpc.server: server stopped")
	return nil
}

// GracefulStop stops the gRPC server gracefully.
// It stops accepting new connections and blocks until all connections are processed
func GracefulStop(ctx context.Context) error {
	log.Info("grpc.server: stopping server")
	if cmuxserver == nil {
		log.Fatalf("grpc.server: server is not listening call ListendAndServe first")
	}

	err := cmuxserver.Shutdown(ctx)
	if err != nil {
		log.WithError(errors.FromError(err).ExtendComponent(component)).
			Errorf("grpc.server: error while gracefully stoping server")
		return err
	}
	return nil
}
