package grpcserver

import (
	"context"
	"net"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/key"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tracing/opentracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const component = "grpc.server"

var (
	initOnce  = &sync.Once{}
	server    *grpc.Server
	enhancers []Enhancer
)

// Init initialize global gRPC server
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if server == nil {
			// Initialize OpenTracing tracer
			opentracing.Init(ctx)

			// Initialize Authentication Manager
			authjwt.Init(ctx)
			authkey.Init(ctx)
			auth := authentication.CombineAuth(authkey.GlobalAuth(), authjwt.GlobalAuth())

			// Declare server with interceptors
			server = NewServer(auth)

			// Apply enhancers on server
			ApplyEnhancers(server, enhancers...)

			// Register server for prometheus metrics
			grpc_prometheus.Register(server)

			// Replace internal gRPC logger with a logrus logger
			grpclog.SetLoggerV2(
				&LogEntry{
					Entry: log.WithFields(log.Fields{"system": "grpc.internal"}),
				},
			)
		}

		// Log registered services
		var services []string
		for name := range server.GetServiceInfo() {
			services = append(services, name)
		}

		log.WithFields(log.Fields{
			"grpc.services": services,
		}).Infof("%s: ready", component)
	})
}

// GlobalServer return global gRPC server
func GlobalServer() *grpc.Server {
	return server
}

// SetGlobalServer sets global gRPC server
func SetGlobalServer(s *grpc.Server) {
	server = s
}

// AddEnhancers adds gRPC server enhancers that will be called at Init time
// Note that it should be called before Init()
func AddEnhancers(fns ...Enhancer) {
	enhancers = append(enhancers, fns...)
}

// ListenAndServe starts global server
func ListenAndServe() {
	if server == nil {
		log.Fatalf("%s: server is not initialized", component)
	}

	lis, err := net.Listen("tcp", URL())
	if err != nil {
		log.WithError(errors.ConnectionError(err.Error()).ExtendComponent(component)).
			WithFields(log.Fields{"grpc.url": URL()}).
			Error("failed to listen")
		return
	}
	log.Infof("%s: start serving on %q", URL(), component)

	// Serve requests
	err = server.Serve(lis)
	if err != nil {
		log.WithError(errors.FromError(err).ExtendComponent(component)).
			WithFields(log.Fields{"grpc.hostname": URL()}).
			Errorf("%s: error listening tcp connections", component)
	} else {
		log.Infof("%s: server gracefully stopped", component)
	}
}

// GracefulStop stops the gRPC server gracefully.
// It stops accepting new connections and blocks until all connections are processed
func StopServer(_ context.Context) {
	server.GracefulStop()
}
