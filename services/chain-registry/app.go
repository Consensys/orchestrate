package chainregistry

import (
	"context"
	"reflect"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/dgraph-io/ristretto"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/httpcache"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/ratelimit"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	chainUCs "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/use-cases/chains"
	faucetsUCs "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/use-cases/faucets"
	ctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers"
	chainctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers/chains"
	faucetctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers/faucets"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

func New(
	cfg *Config,
	pgmngr postgres.Manager,
	ec ethclient.Client,
	jwt, key auth.Checker,
) (*app.App, error) {
	//
	storeBuilder := store.NewBuilder(pgmngr)
	dataAgents, err := storeBuilder.Build(context.Background(), cfg.Store)
	if err != nil {
		return nil, err
	}

	getChainsUC := chainUCs.NewGetChains(dataAgents.Chain)
	getChainUC := chainUCs.NewGetChain(dataAgents.Chain)

	// Create HTTP Handler for Chain
	chainCtrl := chainctrl.NewController(
		getChainsUC,
		getChainUC,
		chainUCs.NewRegisterChain(dataAgents.Chain, ec),
		chainUCs.NewDeleteChain(dataAgents.Chain),
		chainUCs.NewUpdateChain(dataAgents.Chain),
	)

	getFaucetsUC := faucetsUCs.NewGetFaucets(dataAgents.Faucet)
	// Create HTTP Handler for Faucet
	faucetCtrl := faucetctrl.NewController(
		getFaucetsUC,
		faucetsUCs.NewGetFaucet(dataAgents.Faucet),
		faucetsUCs.NewRegisterFaucet(dataAgents.Faucet),
		faucetsUCs.NewDeleteFaucet(dataAgents.Faucet),
		faucetsUCs.NewUpdateFaucet(dataAgents.Faucet),
		faucetsUCs.NewFaucetCandidateUseCase(getChainUC, getFaucetsUC, ec),
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

	cache, err := ristretto.NewCache(cfg.Cache)
	if err != nil {
		return nil, err
	}

	// RateLimit Middleware
	rateLimitOpt := app.MiddlewareOpt(
		reflect.TypeOf(&dynamic.RateLimit{}),
		ratelimit.NewBuilder(ratelimit.NewManager(cache)),
	)

	// HTTPCache Middleware
	httpCacheOpt := app.MiddlewareOpt(
		reflect.TypeOf(&dynamic.HTTPCache{}),
		httpcache.NewBuilder(cache, httpCacheRequest, httpCacheResponse),
	)

	// Create appli to expose metrics
	appli, err := app.New(
		cfg.App,
		app.MultiTenancyOpt("auth", jwt, key, cfg.Multitenancy),
		app.MetricsOpt(),
		app.LoggerMiddlewareOpt("base"),
		rateLimitOpt,
		app.SwaggerOpt("./public/swagger-specs/services/chain-registry/swagger.json", "base@logger-base"),
		chainHandlerOpt,
		httpCacheOpt,
		reverseProxyOpt,
		app.ProviderOpt(
			NewProvider(getChainsUC, time.Second, cfg.ProxyCacheTTL),
		),
	)
	if err != nil {
		return nil, err
	}

	// TODO: chain import should append after starting App not at app creation
	// (or should be deprecated)
	b := func() backoff.BackOff {
		cfg := utils.NewConfig(viper.GetViper())
		cfg.Retry.MaxElapsedTime = 30 * time.Second
		return utils.NewBackOff(cfg)
	}

	httpClient := http.NewClient(http.NewDefaultConfig())
	importChainUC := chainUCs.NewImportChain(dataAgents.Chain, rpc.NewClient(b, httpClient))
	for _, jsonChain := range cfg.EnvChains {
		err = importChainUC.Execute(context.Background(), jsonChain)
		if err != nil {
			if errors.IsAlreadyExistsError(err) {
				appli.Logger().WithError(err).Warnf("skipping import (chain already exists)")
			} else {
				appli.Logger().WithError(err).Errorf("could not import chain")
				return nil, err
			}
		}
	}

	return appli, nil
}
