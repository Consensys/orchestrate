package transactionscheduler

import (
	"context"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/controllers"

	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	pkghttp "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"
)

func newApplication(
	ctx context.Context,
	cfg *Config,
	jwt, key auth.Checker,
	logger *logrus.Logger,
	chainRegistryClient client.ChainRegistryClient,
) (*app.App, error) {
	// TODO: Do all dependency injection in container.go
	// Create Data agents
	pgmngr := postgres.GetManager()
	dataAgents, err := store.Build(ctx, cfg.store, pgmngr)
	if err != nil {
		logger.WithError(err).Fatalf("could not create data-agents")
		return nil, err
	}

	// Create HTTP Router builder
	vals := validators.NewValidators(chainRegistryClient)
	ucs := usecases.NewUseCases(dataAgents, vals)
	ctrls := controllers.NewBuilder(ucs)
	routerBuilder, err := service.NewHTTPBuilder(cfg.app.HTTP, jwt, key, cfg.multitenancy, ctrls)
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
