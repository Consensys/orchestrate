package transactionscheduler

import (
	"context"
	"reflect"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/controllers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"
)

func New(
	cfg *Config,
	pgmngr postgres.Manager,
	jwt, key auth.Checker,
	chainRegistryClient client.ChainRegistryClient,
) (*app.App, error) {
	// Create Data agents
	storeBuilder := store.NewBuilder(pgmngr)
	dataAgents, err := storeBuilder.Build(context.Background(), cfg.Store)
	if err != nil {
		return nil, err
	}

	// Option for transaction handler
	txSchedulerHandlerOpt := app.HandlerOpt(
		reflect.TypeOf(&dynamic.Transactions{}),
		controllers.NewBuilder(usecases.NewUseCases(dataAgents, chainRegistryClient)),
	)

	// Create app
	return app.New(
		cfg.App,
		app.MultiTenancyOpt("auth", jwt, key, cfg.Multitenancy),
		app.MetricsOpt(),
		app.LoggerMiddlewareOpt("base"),
		app.SwaggerOpt("./public/swagger-specs/types/transaction-scheduler/swagger.json", "base@logger-base"),
		txSchedulerHandlerOpt,
		app.ProviderOpt(NewProvider()),
	)
}
