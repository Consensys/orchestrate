package envelopestore

import (
	"context"

	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	pkggrpc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	pkghttp "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/http"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

func newApplication(
	ctx context.Context,
	cfg *Config,
	jwt, key auth.Checker,
	srv svc.EnvelopeStoreServer,
	logger *logrus.Logger,
) (*app.App, error) {
	// Create GRPC builder
	checker := auth.CombineCheckers(key, jwt)
	grpcBuilder, err := grpc.NewServerBuilder(srv, checker, cfg.multitenancy, logger)
	if err != nil {
		return nil, err
	}

	
	// Create GRPC entrypoint
	grpcServer, err := grpcBuilder.BuildServer(ctx, "", grpc.NewStaticConfig())
	if err != nil {
		return nil, err
	}
	grpcEp := pkggrpc.NewEntryPoint(cfg.app.GRPC, grpcServer)

	// Create HTTP Router builder
	routerBuilder, err := http.NewRouterBuilder(srv, cfg.app.HTTP, jwt, key, cfg.multitenancy)
	if err != nil {
		return nil, err
	}

	// Create HTTP EntryPoints
	httpEps := pkghttp.NewEntryPoints(
		cfg.app.HTTP.EntryPoints,
		routerBuilder,
	)

	watcherCfg := configwatcher.NewConfig(cfg.app.HTTP, cfg.app.Watcher)
	watcher := configwatcher.NewWatcher(watcherCfg, httpEps)

	// Create app
	return app.New(watcher, httpEps, grpcEp), nil
}
