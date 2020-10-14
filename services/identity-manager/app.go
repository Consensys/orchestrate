package identitymanager

import (
	"context"
	"reflect"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/builder"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/client"

	"github.com/go-pg/pg/v9/orm"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/service/controllers"
	store "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/store/multi"
)

func NewIdentityManager(cfg *Config, pgmngr postgres.Manager, jwt, key auth.Checker, clt client.KeyManagerClient) (*app.App, error) {
	// Create Data agents
	db, err := store.Build(context.Background(), cfg.Store, pgmngr)
	if err != nil {
		return nil, err
	}

	ucs := builder.NewUseCases(db, clt)

	// Option for identity manager handler
	identityManagerHandlerOpt := app.HandlerOpt(reflect.TypeOf(&dynamic.Identity{}), controllers.NewBuilder(ucs))

	// Create app
	return app.New(
		cfg.App,
		app.MultiTenancyOpt("auth", jwt, key, cfg.Multitenancy),
		ReadinessOpt(db, clt),
		app.MetricsOpt(),
		app.LoggerMiddlewareOpt("base"),
		app.SwaggerOpt("./public/swagger-specs/services/identity-manager/swagger.json", "base@logger-base"),
		identityManagerHandlerOpt,
		app.ProviderOpt(NewProvider()),
	)
}

func ReadinessOpt(db database.DB, clt client.KeyManagerClient) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("database", postgres.Checker(db.(orm.DB)))
		ap.AddReadinessCheck("key-manager", clt.Checker())
		return nil
	}
}
