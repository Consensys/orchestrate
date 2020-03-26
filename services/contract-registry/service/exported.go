package servicelayer

import (
	"context"
	"net/http"
	"path"

	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/contract-registry"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/rest"
)

// InitGRPC Initialize gRPC server
func InitGRPC(cancelCtx context.Context, cancel context.CancelFunc, contractRegistryHandler svc.ContractRegistryServer) {
	grpcserver.AddEnhancers(
		func(s *grpc.Server) *grpc.Server {
			svc.RegisterContractRegistryServer(s, contractRegistryHandler)
			return s
		})
	grpcserver.Init(cancelCtx)
	grpcserver.ListenAndServe()

	cancel()
}

// InitHTTP Initialize REST server
func InitHTTP(cancelCtx context.Context) {
	rest.AddEnhancers(
		func(cancelCtx context.Context, _ *http.ServeMux, gwMux *runtime.ServeMux, conn *grpc.ClientConn) error {
			return svc.RegisterContractRegistryHandler(cancelCtx, gwMux, conn)
		},
		func(ctx context.Context, mux *http.ServeMux, _ *runtime.ServeMux, _ *grpc.ClientConn) error {
			mux.HandleFunc("/swagger/swagger.json", rest.ServeFile(path.Join(rest.SwaggerSpecsPath, "types/contract-registry/registry.swagger.json")))
			return nil
		})
	rest.Init(cancelCtx)
}

// ListenAndServe Serves GRPC and HTTP servers
func ListenAndServe(cancel context.CancelFunc) {
	rest.ListenAndServe()
	cancel()
}
