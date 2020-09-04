package app

import (
	"context"
	"encoding/json"
	"sync"

	traefiklog "github.com/containous/traefik/v2/pkg/log"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/aggregator"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	grpcinterceptor "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/static"
	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server"
	grpcstaticserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server/static"
	grpcservice "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/service/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	httphandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/dynamic"
	httpmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/dynamic"
	metricsmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router"
	httprouter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/multi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tcp"
)

type Daemon interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

type Service interface {
	Start(ctx context.Context) chan error
	Stop(ctx context.Context) error
	IsReady() bool
}

type App struct {
	cfg *Config

	http        *http.EntryPoints
	httpBuilder router.Builder

	grpc        *grpc.EntryPoint
	grpcBuilder grpcserver.Builder

	provider *aggregator.Provider
	watcher  configwatcher.Watcher

	metrics metrics.Registry

	logger *logrus.Logger

	cancel func()

	wg *sync.WaitGroup

	isReady bool
}

func New(cfg *Config, opts ...Option) (*App, error) {
	// Create and configure logger
	logger := logrus.StandardLogger()
	err := log.ConfigureLogger(cfg.Log, logger)
	if err != nil {
		return nil, err
	}

	reg := multi.New(cfg.Metrics)

	httpBuilder := httprouter.NewBuilder(cfg.HTTP.TraefikStatic(), nil)
	httpBuilder.Handler = httphandler.NewBuilder()
	httpBuilder.Middleware = httpmid.NewBuilder()
	httpBuilder.Metrics = metricsmid.NewBuilder(reg.HTTP())

	grpcBuilder := grpcstaticserver.NewBuilder()
	grpcBuilder.Interceptor = grpcinterceptor.NewBuilder()
	grpcBuilder.Service = grpcservice.NewBuilder()

	// Create watcher
	prvdr := aggregator.New()
	watcher := configwatcher.New(
		cfg.Watcher,
		prvdr,
		dynamic.Merge,
		nil,
	)

	// Create App and set provider
	app := newApp(
		cfg,
		httpBuilder,
		grpcBuilder,
		watcher,
		reg,
		logger,
	)
	app.provider = prvdr

	for _, opt := range opts {
		err := opt(app)
		if err != nil {
			return nil, err
		}
	}

	return app, nil
}

func newApp(
	cfg *Config,
	httpBuilder router.Builder,
	grpcBuilder grpcserver.Builder,
	watcher configwatcher.Watcher,
	reg metrics.Registry,
	logger *logrus.Logger,
) *App {
	return &App{
		cfg:         cfg,
		wg:          &sync.WaitGroup{},
		httpBuilder: httpBuilder,
		grpcBuilder: grpcBuilder,
		watcher:     watcher,
		metrics:     reg,
		logger:      logger,
	}
}

func (app *App) init(ctx context.Context) error {
	conf, err := json.Marshal(app.cfg)
	if err != nil {
		return err
	}
	traefiklog.FromContext(ctx).Infof("loaded app configuration %s", string(conf))

	// Create HTTP EntryPoints
	if app.cfg.HTTP != nil {
		app.http = http.NewEntryPoints(app.cfg.HTTP.EntryPoints, app.httpBuilder, app.metrics.TCP())

		// Add Listeners for HTTP
		app.watcher.AddListener(app.http.Switch)
		app.watcher.AddListener(
			func(_ context.Context, cfg interface{}) error {
				if dynCfg, ok := cfg.(*dynamic.Configuration); ok {
					return app.metrics.HTTP().Switch(dynCfg)
				}
				return nil
			},
		)
	}

	// Create GRPC EntryPoint
	if app.cfg.GRPC != nil && app.cfg.GRPC.Static != nil {
		app.grpc = grpc.NewEntryPoint("", app.cfg.GRPC.EntryPoint, app.grpcBuilder, app.metrics.TCP())
		err := app.grpc.BuildServer(ctx, app.cfg.GRPC.Static)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) HTTP() router.Builder {
	return app.httpBuilder
}

func (app *App) GRPC() grpcserver.Builder {
	return app.grpcBuilder
}

func (app *App) Metrics() metrics.Registry {
	return app.metrics
}

func (app *App) Logger() *logrus.Logger {
	return app.logger
}

func (app *App) AddProvider(prvdr provider.Provider) {
	app.provider.AddProvider(prvdr)
}

func (app *App) AddListener(listener func(context.Context, interface{}) error) {
	app.watcher.AddListener(listener)
}

func (app *App) Start(ctx context.Context) error {
	err := app.init(ctx)
	if err != nil {
		return err
	}

	traefiklog.FromContext(ctx).Infof("starting app...")

	cancelableCtx, cancel := context.WithCancel(ctx)
	app.cancel = cancel

	app.wg.Add(3)
	go func() {
		if app.watcher != nil {
			_ = app.watcher.Run(cancelableCtx)
			_ = app.watcher.Close()
		}
		app.wg.Done()
	}()
	go func() {
		if app.http != nil {
			_ = app.http.ListenAndServe(ctx)
		}
		app.wg.Done()
	}()

	go func() {
		if app.grpc != nil {
			_ = app.grpc.ListenAndServe(ctx)
		}
		app.wg.Done()
	}()

	app.isReady = true
	return nil
}

func (app *App) Stop(ctx context.Context) error {
	traefiklog.FromContext(ctx).Infof("gracefully shutting down application...")
	app.cancel()

	app.wg.Add(2)
	var errHTTP error
	go func() {
		if app.http != nil {
			errHTTP = errors.CombineErrors(tcp.Shutdown(ctx, app.http), tcp.Close(app.http))
		}
		app.wg.Done()
	}()

	var errGRPC error
	go func() {
		if app.grpc != nil {
			errGRPC = errors.CombineErrors(tcp.Shutdown(ctx, app.grpc), tcp.Close(app.grpc))
		}
		app.wg.Done()
	}()

	app.wg.Wait()

	if err := errors.CombineErrors(errHTTP, errGRPC); err != nil {
		traefiklog.FromContext(ctx).WithError(err).Errorf("application did not shut down gracefully")
		return err // timed out
	}

	traefiklog.FromContext(ctx).Infof("gracefully shutted down application")
	return nil // completed normally
}

func (app *App) IsReady() bool {
	return app.isReady
}
