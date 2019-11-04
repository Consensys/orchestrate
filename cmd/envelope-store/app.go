package envelopestore

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/healthcheck"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/services/envelope-store"
	envelopestore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store"
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
		envelopestore.Init()

		// Initialize gRPC server
		grpcserver.AddEnhancers(
			func(s *grpc.Server) *grpc.Server {
				svc.RegisterEnvelopeStoreServer(s, envelopestore.GlobalEnvelopeStoreServer())
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
