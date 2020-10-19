package identitymanager

import (
	"context"
	"reflect"

	client2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/builder"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"
	client3 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"

	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/service/controllers"
	store "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/multi"
)

func NewIdentityManager(cfg *Config, pgmngr postgres.Manager, jwt, key auth.Checker, keyManagerClient client.KeyManagerClient,
	registryClient client2.ChainRegistryClient, txSchedulerClient client3.TransactionSchedulerClient) (*app.App, error) {
	// Create Data agents
	db, err := store.Build(context.Background(), cfg.Store, pgmngr)
	if err != nil {
		return nil, err
	}

	ucs := builder.NewUseCases(db, keyManagerClient, registryClient, txSchedulerClient)

	// Option for identity manager handler
	identityManagerHandlerOpt := app.HandlerOpt(reflect.TypeOf(&dynamic.Identity{}), controllers.NewBuilder(ucs))

	// Create app
	return app.New(
		cfg.App,
		app.MultiTenancyOpt("auth", jwt, key, cfg.Multitenancy),
		ReadinessOpt(db, keyManagerClient, registryClient, txSchedulerClient),
		app.MetricsOpt(),
		app.LoggerMiddlewareOpt("base"),
		app.SwaggerOpt("./public/swagger-specs/services/identity-manager/swagger.json", "base@logger-base"),
		identityManagerHandlerOpt,
		app.ProviderOpt(NewProvider()),
	)
}

func ReadinessOpt(db database.DB, keyManagerClient client.KeyManagerClient, registryClient client2.ChainRegistryClient,
	txSchedulerClient client3.TransactionSchedulerClient) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("database", postgres.Checker(db.(orm.DB)))
		ap.AddReadinessCheck("key-manager", keyManagerClient.Checker())
		ap.AddReadinessCheck("transaction-scheduler", txSchedulerClient.Checker())
		ap.AddReadinessCheck("chain-registry", registryClient.Checker())
		return nil
	}
}
