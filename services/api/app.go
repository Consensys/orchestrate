package api

import (
	"context"
	"reflect"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/middleware/httpcache"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/middleware/ratelimit"
	"github.com/ConsenSys/orchestrate/services/api/proxy"
	"github.com/dgraph-io/ristretto"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/ethclient"

	qkmclient "github.com/consensys/quorum-key-manager/pkg/client"

	pkgsarama "github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	pkgproxy "github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/handler/proxy"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/database"
	"github.com/ConsenSys/orchestrate/services/api/business/builder"
	"github.com/ConsenSys/orchestrate/services/api/metrics"
	"github.com/Shopify/sarama"
	"github.com/go-pg/pg/v9/orm"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/auth"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/ConsenSys/orchestrate/services/api/service/controllers"
	"github.com/ConsenSys/orchestrate/services/api/store/multi"
)

func NewAPI(
	cfg *Config,
	pgmngr postgres.Manager,
	jwt, key auth.Checker,
	keyManagerClient qkmclient.Eth1Client,
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
