package chainregistry

import (
	"context"
	"reflect"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/go-pg/pg/v9"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth"
	pkgpg "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/middleware/httpcache"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/middleware/ratelimit"
	chainUCs "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/chain-registry/use-cases/chains"
	ctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/service/controllers"
	chainctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/service/controllers/chains"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/multi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/postgres"
)

func New(
	cfg *Config,
	pgmngr pkgpg.Manager,
	ec ethclient.Client,
	jwt, key auth.Checker,
) (*app.App, error) {
	db, err := multi.Build(context.Background(), cfg.Store, pgmngr)
	if err != nil {
		return nil, err
	}

	dataAgents := postgres.Build(db)

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

	chainHandlerOpt := app.HandlerOpt(
		reflect.TypeOf(&dynamic.Chains{}),
		ctrl.NewBuilder(chainCtrl),
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
		ReadinessOpt(db),
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

	importChainUC := chainUCs.NewImportChain(dataAgents.Chain, rpc.GlobalClient())
	for _, jsonChain := range cfg.EnvChains {
		err = importChainUC.Execute(utils.RetryConnectionError(context.Background(), true), jsonChain)
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

func ReadinessOpt(db *pg.DB) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("database", pkgpg.Checker(db))
		return nil
	}
}
