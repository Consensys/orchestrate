package keymanager

import (
	"context"
	"reflect"

	"github.com/ConsenSys/orchestrate/services/key-manager/key-manager/builder"

	"github.com/ConsenSys/orchestrate/services/key-manager/store"

	"github.com/ConsenSys/orchestrate/pkg/app"
	"github.com/ConsenSys/orchestrate/pkg/http/config/dynamic"
	"github.com/ConsenSys/orchestrate/services/key-manager/service/controllers"
)

func NewKeyManager(ctx context.Context, cfg *Config) (*app.App, error) {
	// Create Data agents
	vault, err := store.Build(ctx, cfg.Store)
	if err != nil {
		return nil, err
	}

	// Option for key manager handler
	keyManagerHandlerOpt := app.HandlerOpt(reflect.TypeOf(&dynamic.KeyManager{}),
		controllers.NewBuilder(vault, builder.NewETHUseCases(vault), builder.NewZKSUseCases(vault)))

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
