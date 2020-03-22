package envelopestore

import (
	"context"

	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store"
)

func New(
	cfg *app.Config,
	jwt, key auth.Checker,
	multitenancy bool,
	service svc.EnvelopeStoreServer,
	logger *logrus.Logger,
) (*app.App, error) {
	// Create GRPC builder
	checker := auth.CombineCheckers(key, jwt)
	grpcBuilder, err := NewGRPCBuilder(service, checker, multitenancy, logger)
	if err != nil {
		return nil, err
	}

	// Create GRPC entrypoint
	server, err := grpcBuilder.Build(
		context.Background(),
		"",
		NewGRPCStaticConfig(),
	)
	if err != nil {
		return nil, err
	}
	grpcEp := grpc.NewEntryPoint(cfg.GRPC, server)

	// Create HTTP Router builder
	httpBuilder, err := NewHTTPBuilder(cfg.HTTP, jwt, key, multitenancy, service)
	if err != nil {
		return nil, err
	}

	// Create HTTP EntryPoints
	httpEps := http.NewEntryPoints(
		cfg.HTTP.EntryPoints,
		httpBuilder,
	)

	// Create Configuration Watcher
	// Create configuration listener switching HTTP Endpoints configuration
	listeners := []func(context.Context, interface{}) error{
		httpEps.Switch,
	}

	watcher := configwatcher.New(
		cfg.Watcher,
		NewProvider(cfg.HTTP),
		dynamic.Merge,
		listeners,
	)

	// Create app
	return app.New(watcher, httpEps, grpcEp), nil
}

func NewService(pgmngr postgres.Manager, storeCfg *store.Config) (svc.EnvelopeStoreServer, error) {
	// Create Store
	storeBuilder := store.NewBuilder(pgmngr)
	return storeBuilder.Build(context.Background(), storeCfg)
}
