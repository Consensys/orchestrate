package transactionscheduler

import (
	"context"

	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	pkghttp "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"
)

func newApplication(
	ctx context.Context,
	cfg *Config,
	jwt, key auth.Checker,
	logger *logrus.Logger,
) (*app.App, error) {
	// Create Data agents
	pgmngr := postgres.GetManager()
	dataAgents, err := store.Build(ctx, cfg.store, pgmngr)
	if err != nil {
		logger.WithError(err).Fatalf("could not create data-agents")
		return nil, err
	}

	// Create HTTP Router builder
	routerBuilder, err := service.NewHTTPBuilder(cfg.app.HTTP, jwt, key, cfg.multitenancy, usecases.NewUseCases(dataAgents))
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
	return app.New(watcher, httpEps, nil), nil
}
