package contractregistry

import (
	"context"
	"net/http"
	"path"
	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/rest"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"
)

var (
	app       = common.NewApp()
	startOnce = &sync.Once{}
)

// Run application
func Start(ctx context.Context) {
	startOnce.Do(func() {
		// Initialize gRPC Server service
		contractregistry.Init(ctx)

		cancelCtx, cancel := context.WithCancel(ctx)
		go metrics.StartServer(ctx, cancel, app.IsAlive, app.IsReady)

		go func() {
			// Initialize gRPC server
			grpcserver.AddEnhancers(
				func(s *grpc.Server) *grpc.Server {
					types.RegisterRegistryServer(s, contractregistry.GlobalRegistry())
					return s
				})
			grpcserver.Init(cancelCtx)
			grpcserver.ListenAndServe()
			cancel()
		}()

		// Initialize REST server
		rest.AddEnhancers(
			func(cancelCtx context.Context, _ *http.ServeMux, gwMux *runtime.ServeMux, conn *grpc.ClientConn) error {
				return types.RegisterRegistryHandler(cancelCtx, gwMux, conn)
			},
			func(ctx context.Context, mux *http.ServeMux, _ *runtime.ServeMux, _ *grpc.ClientConn) error {
				mux.HandleFunc("/swagger/swagger.json", rest.ServeFile(path.Join(rest.SwaggerSpecsPath, "types/contract-registry/registry.swagger.json")))
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
	app.SetReady(false)
	common.InParallel(
		func() { grpcserver.StopServer(ctx) },
		func() { metrics.StopServer(ctx) },
		func() { rest.StopServer(ctx) },
	)
	log.Info("app: gracefully stopped application")
}
