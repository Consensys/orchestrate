package chainregistry

import (
	"context"
	"reflect"
	"time"

	"github.com/dgraph-io/ristretto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/ratelimit"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/use-cases"
	ctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers"
	chainctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers/chains"
	faucetctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers/faucets"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

func New(
	cfg *Config,
	pgmngr postgres.Manager,
	ec ethclient.ChainLedgerReader,
	jwt, key auth.Checker,
) (*app.App, error) {
	//
	storeBuilder := store.NewBuilder(pgmngr)
	dataAgents, err := storeBuilder.Build(context.Background(), cfg.Store)
	if err != nil {
		return nil, err
	}

	// Create HTTP Handler for Chain
	chainCtrl := chainctrl.NewController(
		usecases.NewGetChains(dataAgents.Chain),
		usecases.NewGetChain(dataAgents.Chain),
		usecases.NewRegisterChain(dataAgents.Chain, ec),
		usecases.NewDeleteChain(dataAgents.Chain),
		usecases.NewUpdateChain(dataAgents.Chain),
	)

	// Create HTTP Handler for Faucet
	faucetCtrl := faucetctrl.NewController(
		usecases.NewGetFaucets(dataAgents.Faucet),
		usecases.NewGetFaucet(dataAgents.Faucet),
		usecases.NewRegisterFaucet(dataAgents.Faucet),
		usecases.NewDeleteFaucet(dataAgents.Faucet),
		usecases.NewUpdateFaucet(dataAgents.Faucet),
	)
	chainHandlerOpt := app.HandlerOpt(
		reflect.TypeOf(&dynamic.Chains{}),
		ctrl.NewBuilder(chainCtrl, faucetCtrl),
	)

	// ReverseProxy Handler
	proxyBuilder, err := proxy.NewBuilder(cfg.ServersTransport, nil)
	if err != nil {
		return nil, err
	}
	reverseProxyOpt := app.HandlerOpt(
		reflect.TypeOf(&dynamic.ReverseProxy{}),
		proxyBuilder,
	)

	// RateLimit Middleware
	cache, err := ristretto.NewCache(cfg.Cache)
	if err != nil {
		return nil, err
	}
	rateLimitOpt := app.MiddlewareOpt(
		reflect.TypeOf(&dynamic.RateLimit{}),
		ratelimit.NewBuilder(ratelimit.NewManager(cache)),
	)

	// Create appli to expose metrics
	appli, err := app.New(
		cfg.App,
		app.MultiTenancyOpt("auth", jwt, key, cfg.Multitenancy),
		app.MetricsOpt(),
		app.LoggerMiddlewareOpt("base"),
		rateLimitOpt,
		app.SwaggerOpt("./public/swagger-specs/types/chain-registry/swagger.json", "base@logger-base"),
		chainHandlerOpt,
		reverseProxyOpt,
		app.ProviderOpt(
			NewProvider(usecases.NewGetChains(dataAgents.Chain), time.Second),
		),
	)
	if err != nil {
		return nil, err
	}

	// TODO: chain import should append after starting App not at app creation
	// (or should be deprecated)
	importChainUC := usecases.NewImportChain(dataAgents.Chain, ec)
	for _, jsonChain := range cfg.EnvChains {
		err = importChainUC.Execute(context.Background(), jsonChain)
		if err != nil {
			appli.Logger().WithError(err).Errorf("could not import chain")
		}
	}

	return appli, nil
}
