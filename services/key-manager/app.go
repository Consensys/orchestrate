package keymanager

import (
	"context"
	"reflect"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/key-manager/builder"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/service/controllers"
	multistore "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/key-manager/store/multi"
)

func NewTxSigner(cfg *Config) (*app.App, error) {
	// Create Data agents
	vault, err := multistore.Build(context.Background(), cfg.Store)
	if err != nil {
		return nil, err
	}

	// Option for identity manager handler
	txSignerHandlerOpt := app.HandlerOpt(reflect.TypeOf(&dynamic.Signer{}), controllers.NewBuilder(builder.NewUseCases(vault)))

	// Create app
	return app.New(
		cfg.App,
		ReadinessOpt(vault),
		app.MetricsOpt(),
		app.LoggerMiddlewareOpt("base"),
		app.SwaggerOpt("./public/swagger-specs/services/tx-signer/swagger.json", "base@logger-base"),
		txSignerHandlerOpt,
		app.ProviderOpt(NewProvider()),
	)
}

func ReadinessOpt(vault store.Vault) app.Option {
	return func(ap *app.App) error {
		// TODO: Add readiness check for the Vault
		return nil
	}
}
