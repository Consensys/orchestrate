package keymanager

import (
	"context"
	"reflect"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/service/controllers"
)

func NewKeyManager(ctx context.Context, cfg *Config) (*app.App, error) {
	// Create Data agents
	vault, err := store.Build(ctx, cfg.Store)
	if err != nil {
		return nil, err
	}

	// Option for key manager handler
	keyManagerHandlerOpt := app.HandlerOpt(reflect.TypeOf(&dynamic.Signer{}), controllers.NewBuilder(vault))

	// Create app
	return app.New(
		cfg.App,
		ReadinessOpt(vault),
		app.MetricsOpt(),
		app.LoggerMiddlewareOpt("base"),
		app.SwaggerOpt("./public/swagger-specs/services/key-manager/swagger.json", "base@logger-base"),
		keyManagerHandlerOpt,
		app.ProviderOpt(NewProvider()),
	)
}

func ReadinessOpt(vault store.Vault) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("vault", vault.HealthCheck)
		return nil
	}
}
