package keymanager

import (
	"context"
	"reflect"

	healthz "github.com/heptiolabs/healthcheck"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/key-manager/builder"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/service/controllers"
	multistore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store/multi"
)

func NewKeyManager(ctx context.Context, cfg *Config) (*app.App, error) {
	// Create Data agents
	vault, err := multistore.Build(ctx, cfg.Store)
	if err != nil {
		return nil, err
	}

	// Create UCs
	ethUCs := builder.NewEthereumUseCases(vault)

	// Option for key manager handler
	keyManagerHandlerOpt := app.HandlerOpt(reflect.TypeOf(&dynamic.Signer{}), controllers.NewBuilder(ethUCs))

	// Create app
	return app.New(
		cfg.App,
		ReadinessOpt(vault.HealthCheck()),
		app.MetricsOpt(),
		app.LoggerMiddlewareOpt("base"),
		app.SwaggerOpt("./public/swagger-specs/services/key-manager/swagger.json", "base@logger-base"),
		keyManagerHandlerOpt,
		app.ProviderOpt(NewProvider()),
	)
}

func ReadinessOpt(checker healthz.Check) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("vault", checker)
		return nil
	}
}
