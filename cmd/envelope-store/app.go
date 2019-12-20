package envelopestore

import (
	"context"
	"net/http"
	"path"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/rest"
	envelopestore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope-store"
	"google.golang.org/grpc"
)

var (
	app       = common.NewApp()
	startOnce = &sync.Once{}
)

// Run application
func Start(ctx context.Context) {
	startOnce.Do(func() {
		multitenancy.Init(ctx)

		// Initialize gRPC Server service
		envelopestore.Init()

		cancelCtx, cancel := context.WithCancel(ctx)
		go metrics.StartServer(ctx, cancel, app.IsAlive, app.IsReady)

		go func() {
			// Initialize gRPC server
			grpcserver.AddEnhancers(
				func(s *grpc.Server) *grpc.Server {
					types.RegisterEnvelopeStoreServer(s, envelopestore.GlobalEnvelopeStoreServer())
					return s
				})
			grpcserver.Init(cancelCtx)
			grpcserver.ListenAndServe()
			cancel()
		}()

		// Initialize REST server
		rest.AddEnhancers(
			func(ctx context.Context, _ *http.ServeMux, gwMux *runtime.ServeMux, conn *grpc.ClientConn) error {
				return types.RegisterEnvelopeStoreHandler(cancelCtx, gwMux, conn)
			},
			func(ctx context.Context, mux *http.ServeMux, _ *runtime.ServeMux, _ *grpc.ClientConn) error {
				mux.HandleFunc("/swagger/swagger.json", rest.ServeFile(path.Join(rest.SwaggerSpecsPath, "types/envelope-store/store.swagger.json")))
				return nil
			})
		rest.Init(cancelCtx)

		// Indicate that application is ready
		app.SetReady(true)

		rest.ListenAndServe()
		cancel()
	})
}

// Close gracefully stops the application
func Close(ctx context.Context) {
	log.Warn("app: stopping...")
	common.InParallel(
		func() { grpcserver.StopServer(ctx) },
		func() { metrics.StopServer(ctx) },
		func() { rest.StopServer(ctx) },
	)
	log.Info("app: gracefully stopped application")
}
