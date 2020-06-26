package envelopestore

import (
	"context"
	"reflect"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	grpcstatic "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/service/controllers"
	grpcservice "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/service/grpc"
	httpservice "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/service/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/store"
)

func New(
	cfg *Config,
	pgmngr postgres.Manager,
	jwt, key auth.Checker,
) (*app.App, error) {
	// Create Store
	storeBuilder := store.NewBuilder(pgmngr)
	dataAgents, err := storeBuilder.Build(context.Background(), cfg.Store)
	if err != nil {
		return nil, err
	}

	srv, err := controllers.NewGRPCService(dataAgents)
	if err != nil {
		return nil, err
	}

	envelopeServiceOpt := app.ServiceOpt(
		reflect.TypeOf(&grpcstatic.Envelopes{}),
		grpcservice.NewBuilder(srv),
	)

	envelopeHandlerOpt := app.HandlerOpt(
		reflect.TypeOf(&dynamic.Envelopes{}),
		httpservice.NewBuilder(srv),
	)

	// Create app
	return app.New(
		cfg.App,
		app.MultiTenancyOpt("auth", jwt, key, cfg.Multitenancy),
		app.MetricsOpt(),
		app.LoggerMiddlewareOpt("base"),
		app.LoggerInterceptorOpt(),
		app.SwaggerOpt("./public/swagger-specs/services/envelope-store/proto/store.swagger.json", "base@logger-base"),
		envelopeServiceOpt,
		envelopeHandlerOpt,
		app.ProviderOpt(NewProvider()),
	)
}
