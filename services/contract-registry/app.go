package contractregistry

import (
	"context"
	"reflect"

	"github.com/go-pg/pg/v9"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	pkgpg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	grpcstatic "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/abi"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/contract-registry/use-cases"
	grpcservice "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/service/grpc"
	httpservice "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/service/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/multi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/postgres"
)

func New(
	cfg *Config,
	pgmngr pkgpg.Manager,
	jwt, key auth.Checker,
) (*app.App, error) {
	db, err := multi.Build(context.Background(), cfg.Store, pgmngr)
	if err != nil {
		return nil, err
	}

	contractDA, repositoryDA, tagDA, artifactDA, methodDA, eventDA, codeHashDA := postgres.Build(db)

	registerContractUC := usecases.NewRegisterContract(contractDA)
	getContractUC := usecases.NewGetContract(artifactDA)

	srv := grpcservice.New(
		registerContractUC,
		getContractUC,
		usecases.NewGetMethods(methodDA),
		usecases.NewGetEvents(eventDA),
		usecases.NewGetCatalog(repositoryDA),
		usecases.NewGetTags(tagDA),
		usecases.NewSetCodeHash(codeHashDA),
		usecases.NewGetMethodSignatures(getContractUC),
	)

	contractServiceOpt := app.ServiceOpt(
		reflect.TypeOf(&grpcstatic.Contracts{}),
		grpcservice.NewBuilder(srv),
	)

	contractHandlerOpt := app.HandlerOpt(
		reflect.TypeOf(&dynamic.Contracts{}),
		httpservice.NewBuilder(srv),
	)

	appli, err := app.New(
		cfg.App,
		app.MultiTenancyOpt("auth", jwt, key, cfg.Multitenancy),
		ReadinessOpt(db),
		app.MetricsOpt(),
		app.LoggerOpt("base"),
		app.SwaggerOpt("./public/swagger-specs/services/contract-registry/proto/registry.swagger.json", "base@logger-base"),
		contractServiceOpt,
		contractHandlerOpt,
		app.ProviderOpt(NewProvider()),
	)
	if err != nil {
		return nil, err
	}

	// TODO: ABI import should append after starting App not at app creation
	// (or should be deprecated)
	for _, ABI := range cfg.ABIs {
		c, err := abi.StringToContract(ABI)
		if err != nil {
			appli.Logger().WithError(err).Errorf("could not parse contract ABI")
			continue
		}

		err = registerContractUC.Execute(context.Background(), c)
		if err != nil {
			appli.Logger().WithError(err).Errorf("could not import contract ABI")
		}
	}

	return appli, nil
}

func ReadinessOpt(db *pg.DB) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("database", pkgpg.Checker(db))
		return nil
	}
}
