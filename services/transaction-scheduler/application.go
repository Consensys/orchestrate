package transactionscheduler

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/controllers"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	metrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/multi"
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
	promRegisterer prom.Registerer,
) (*app.App, error) {
	// Create metrics registry and register it on prometheus
	reg := metrics.New(cfg.app.Metrics)
	err := promRegisterer.Register(reg.Prometheus())

	if err != nil {
		return nil, err
	}

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
	routerBuilder, err := service.NewHTTPBuilder(cfg.app.HTTP, jwt, key, cfg.multitenancy, ctrls, reg.HTTP())
	if err != nil {
		return nil, err
	}

	// Create app
	return app.New(
		cfg.app,
		configwatcher.NewProvider(configwatcher.NewConfig(cfg.app.HTTP, cfg.app.Watcher).DynamicCfg()),
		routerBuilder,
		nil,
		reg,
	)
}
