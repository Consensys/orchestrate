package contractregistry

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	metrics "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/multi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/abi"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/contract-registry/use-cases"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	grpcservice "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/service/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store"
)

func New(
	cfg *app.Config,
	jwt, key auth.Checker,
	multitenancy bool,
	service svc.ContractRegistryServer,
	logger *logrus.Logger,
	registry prom.Registerer,
) (*app.App, error) {
	// Create metrics registry and register it on prometheus
	reg := metrics.New(cfg.Metrics)
	err := registry.Register(reg.Prometheus())
	if err != nil {
		return nil, err
	}

	// Create GRPC builder
	checker := auth.CombineCheckers(key, jwt)
	grpcBuilder, err := NewGRPCBuilder(service, checker, multitenancy, logger, reg.GRPCServer())
	if err != nil {
		return nil, err
	}
	cfg.GRPC.Static = NewGRPCStaticConfig()

	// Create HTTP Router builder
	httpBuilder, err := NewHTTPBuilder(cfg.HTTP, jwt, key, multitenancy, service, reg.HTTP())
	if err != nil {
		return nil, err
	}

	// Create app
	return app.New(
		cfg,
		NewProvider(cfg.HTTP),
		httpBuilder,
		grpcBuilder,
		reg,
	)
}

func NewService(ctx context.Context, pgmngr postgres.Manager, storeCfg *store.Config) (svc.ContractRegistryServer, error) {
	logger := log.FromContext(ctx)
	// Create Store
	storeBuilder := store.NewBuilder(pgmngr)
	contractDA, repositoryDA, tagDA, artifactDA, methodDA, eventDA, codeHashDA, err := storeBuilder.Build(context.Background(), storeCfg)
	if err != nil {
		return nil, err
	}

	registerContractUC := usecases.NewRegisterContract(contractDA)
	err = initializeABIs(ctx, registerContractUC)
	if err != nil {
		logger.WithError(err).Warn("could not initialize contract-registry")
	}

	// Create and return service
	return grpcservice.New(
		registerContractUC,
		usecases.NewGetContract(artifactDA),
		usecases.NewGetMethods(methodDA),
		usecases.NewGetEvents(eventDA),
		usecases.NewGetCatalog(repositoryDA),
		usecases.NewGetTags(tagDA),
		usecases.NewSetCodeHash(codeHashDA),
	), nil
}

// initializeABIs Read ABIs from ABI viper configuration
func initializeABIs(ctx context.Context, registerContractUC usecases.RegisterContractUseCase) error {

	var contracts []*abi.Contract
	for _, ABI := range viper.GetStringSlice(ABIViperKey) {
		c, err := abi.StringToContract(ABI)
		if err != nil {
			return err
		}
		contracts = append(contracts, c)
	}

	// Register contracts
	for _, contract := range contracts {
		err := registerContractUC.Execute(ctx, contract)
		if err != nil {
			return err
		}
	}

	return nil
}
