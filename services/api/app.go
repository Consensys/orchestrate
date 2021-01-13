package api

import (
	"context"
	"reflect"
	"time"

	"github.com/dgraph-io/ristretto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/middleware/httpcache"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/middleware/ratelimit"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/proxy"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"

	keymanager "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"

	"github.com/Shopify/sarama"
	"github.com/go-pg/pg/v9/orm"
	pkgsarama "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database"
	pkgproxy "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/handler/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/builder"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/metrics"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database/postgres"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/service/controllers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/multi"
)

func NewAPI(
	cfg *Config,
	pgmngr postgres.Manager,
	jwt, key auth.Checker,
	keyManagerClient keymanager.KeyManagerClient,
	ec ethclient.Client,
	syncProducer sarama.SyncProducer,
	topicCfg *pkgsarama.KafkaTopicConfig,
) (*app.App, error) {
	// Create Data agents
	db, err := multi.Build(context.Background(), cfg.Store, pgmngr)
	if err != nil {
		return nil, err
	}

	var appMetrics metrics.TransactionSchedulerMetrics
	if cfg.App.Metrics.IsActive(metrics.ModuleName) {
		appMetrics = metrics.NewTransactionSchedulerMetrics()
	} else {
		appMetrics = metrics.NewTransactionSchedulerNopMetrics()
	}

	ucs := builder.NewUseCases(db, appMetrics, keyManagerClient, ec, syncProducer, topicCfg)

	// Option of the API
	apiHandlerOpt := app.HandlerOpt(reflect.TypeOf(&dynamic.API{}), controllers.NewBuilder(ucs, keyManagerClient))

	// ReverseProxy Handler
	proxyBuilder, err := pkgproxy.NewBuilder(cfg.Proxy.ServersTransport, nil)
	if err != nil {
		return nil, err
	}
	reverseProxyOpt := app.HandlerOpt(
		reflect.TypeOf(&dynamic.ReverseProxy{}),
		proxyBuilder,
	)

	cache, err := ristretto.NewCache(cfg.Proxy.Cache)
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
		httpcache.NewBuilder(cache, proxy.HTTPCacheRequest, proxy.HTTPCacheResponse),
	)

	// Create app
	return app.New(
		cfg.App,
		app.MultiTenancyOpt("auth", jwt, key, cfg.Multitenancy),
		ReadinessOpt(db),
		app.MetricsOpt(appMetrics),
		app.LoggerMiddlewareOpt("base"),
		rateLimitOpt,
		app.SwaggerOpt("./public/swagger-specs/services/api/swagger.json", "base@logger-base"),
		apiHandlerOpt,
		httpCacheOpt,
		reverseProxyOpt,
		app.ProviderOpt(NewProvider(ucs.SearchChains(), time.Second, cfg.Proxy.ProxyCacheTTL)),
	)
}

func ReadinessOpt(db database.DB) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("database", postgres.Checker(db.(orm.DB)))
		ap.AddReadinessCheck("kafka", pkgsarama.GlobalClientChecker())
		return nil
	}
}
