package api

import (
	"context"
	"reflect"
	"time"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http/middleware/httpcache"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/middleware/ratelimit"
	ristretto2 "github.com/consensys/orchestrate/pkg/toolkit/cache/ristretto"
	"github.com/consensys/orchestrate/services/api/proxy"
	"github.com/dgraph-io/ristretto"

	"github.com/consensys/orchestrate/pkg/toolkit/ethclient"

	qkmclient "github.com/consensys/quorum-key-manager/pkg/client"

	"github.com/Shopify/sarama"
	pkgsarama "github.com/consensys/orchestrate/pkg/broker/sarama"
	pkgproxy "github.com/consensys/orchestrate/pkg/toolkit/app/http/handler/proxy"
	"github.com/consensys/orchestrate/pkg/toolkit/database"
	"github.com/consensys/orchestrate/services/api/business/builder"
	"github.com/consensys/orchestrate/services/api/metrics"
	"github.com/go-pg/pg/v9/orm"

	"github.com/consensys/orchestrate/pkg/toolkit/app"
	"github.com/consensys/orchestrate/pkg/toolkit/app/auth"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/database/postgres"
	"github.com/consensys/orchestrate/services/api/service/controllers"
	"github.com/consensys/orchestrate/services/api/store/multi"
)

func NewAPI(
	cfg *Config,
	pgmngr postgres.Manager,
	jwt, key auth.Checker,
	keyManagerClient qkmclient.KeyManagerClient,
	qkmStoreID string,
	ec ethclient.Client,
	syncProducer sarama.SyncProducer,
	topicCfg *pkgsarama.KafkaTopicConfig,
) (*app.App, error) {
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

	ucs := builder.NewUseCases(db, appMetrics, keyManagerClient, qkmStoreID, ec, syncProducer, topicCfg)

	// Option of the API
	apiHandlerOpt := app.HandlerOpt(reflect.TypeOf(&dynamic.API{}), controllers.NewBuilder(ucs, keyManagerClient, qkmStoreID))

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
		httpcache.NewBuilder(ristretto2.NewCacheManager(httpcache.Component, cache, cfg.Proxy.ProxyCacheTTL), proxy.HTTPCacheRequest, proxy.HTTPCacheResponse),
	)

	var accessLogMid app.Option
	if cfg.App.HTTP.AccessLog {
		accessLogMid = app.LoggerMiddlewareOpt("base")
	} else {
		accessLogMid = app.NonOpt()
	}

	// Create app
	return app.New(
		cfg.App,
		app.MultiTenancyOpt("auth", jwt, key, cfg.Multitenancy),
		ReadinessOpt(db),
		app.MetricsOpt(appMetrics),
		accessLogMid,
		rateLimitOpt,
		apiHandlerOpt,
		httpCacheOpt,
		reverseProxyOpt,
		app.ProviderOpt(NewProvider(ucs.SearchChains(), time.Second, cfg.Proxy.ProxyCacheTTL, cfg.App.HTTP.AccessLog)),
	)
}

func ReadinessOpt(db database.DB) app.Option {
	return func(ap *app.App) error {
		ap.AddReadinessCheck("database", postgres.Checker(db.(orm.DB)))
		ap.AddReadinessCheck("kafka", pkgsarama.GlobalClientChecker())
		return nil
	}
}
