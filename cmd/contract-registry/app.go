package contractregistry

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/common"
	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/grpc/server"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/http/healthcheck"
	types "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/services/contract-registry"
	services "gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/services/contract-registry"
)

var (
	app       *common.App
	startOnce = &sync.Once{}
)

func init() {
	// Create app
	app = common.NewApp()
}

// Run application
func Start(ctx context.Context) {
	startOnce.Do(func() {
		// Initialize gRPC Server service
		services.Init(ctx)

		// Initialize gRPC server
		grpcserver.AddEnhancers(
			func(s *grpc.Server) *grpc.Server {
				types.RegisterRegistryServer(s, services.GlobalRegistry())
				return s
			},
		)
		grpcserver.Init(ctx)

		// Initialize HTTP server for healthchecks
		http.Init(ctx)
		http.Enhance(healthcheck.HealthCheck(app))

		// Indicate that application is ready
		app.SetReady(true)

		// Start listening
		err := grpcserver.ListenAndServe()
		if err != nil {
			log.WithError(err).Error("app: error listening")
		}
	})
}

// Close gracefully stops the application
func Close(ctx context.Context) {
	log.Warn("app: stopping...")
	err := grpcserver.GracefulStop(ctx)
	if err != nil {
		log.WithError(err).Error("app: error stopping application")
	}
}