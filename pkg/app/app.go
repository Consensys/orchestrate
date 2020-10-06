package app

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	traefiklog "github.com/containous/traefik/v2/pkg/log"
	"github.com/hashicorp/go-multierror"
	healthz "github.com/heptiolabs/healthcheck"
	"github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/aggregator"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	grpcinterceptor "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/interceptor/static"
	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server"
	grpcstaticserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server/static"
	grpcservice "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/service/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	httphandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/healthcheck"
	httpmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/dynamic"
	metricsmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router"
	httprouter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics/multi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tcp"
)

// Daemon are structures exposing a long time running function
// that will be maintained by the App object
type Daemon interface {
	// Run should start a long running session that should stop
	// following a cancel of ctx

	// In case, Run() returns an error, the App automatically
	// triggers a complete Shutdown procedure
	// So a Daemon should do its best to possibly recover
	// before returning an error
	Run(ctx context.Context) error

	// Close allows a daemon to possibly clean its state
	Close() error
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

	daemons        []Daemon
	readinessCheck []*healthcheck.Checker

	logger *logrus.Logger

	cancel func()

	daemonWg sync.WaitGroup
	serverWg sync.WaitGroup

	errors chan error

	closeOnce sync.Once

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
		httpBuilder: httpBuilder,
		grpcBuilder: grpcBuilder,
		watcher:     watcher,
		metrics:     reg,
		logger:      logger,
		errors:      make(chan error),
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

func (app *App) RegisterDaemon(d Daemon) {
	app.daemons = append(app.daemons, d)
}

func (app *App) AddReadinessCheck(name string, check healthz.Check) {
	app.readinessCheck = append(app.readinessCheck, healthcheck.NewChecker(name, check))
}

func (app *App) Start(ctx context.Context) error {
	err := app.init(ctx)
	if err != nil {
		return err
	}

	traefiklog.FromContext(ctx).Infof("starting app...")

	if app.http != nil {
		app.serverWg.Add(1)
		go func() {
			for err := range app.http.ListenAndServe(ctx) {
				if err != nil && err != nethttp.ErrServerClosed {
					app.errors <- err
				}
			}
			app.serverWg.Done()
		}()
	}

	if app.grpc != nil {
		app.serverWg.Add(1)
		go func() {
			err := app.grpc.ListenAndServe(ctx)
			if err != nil && err != nethttp.ErrServerClosed {
				app.errors <- err
			}
			app.serverWg.Done()
		}()
	}

	cancelableCtx, cancel := context.WithCancel(ctx)
	app.cancel = cancel

	if app.watcher != nil {
		app.daemonWg.Add(1)
		go func() {
			err := app.watcher.Run(cancelableCtx)
			if err != nil && err != context.Canceled {
				app.errors <- err
			}
			app.daemonWg.Done()
		}()
	}

	app.daemonWg.Add(len(app.daemons))
	for _, daemon := range app.daemons {
		go func(daemon Daemon) {
			err := daemon.Run(cancelableCtx)
			if err != nil && err != context.Canceled {
				app.errors <- err
			}
			app.daemonWg.Done()
		}(daemon)
	}

	app.isReady = true
	return nil
}

func (app *App) Run(ctx context.Context) error {
	// Start app
	err := app.Start(ctx)
	if err != nil {
		return err
	}

	signals := make(chan os.Signal, 3)
	signal.Notify(signals)

signalLoop:
	for {
		select {
		case sig := <-signals:
			switch sig {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				traefiklog.FromContext(ctx).Infof("signal %q intercepted", sig.String())
				break signalLoop
			case syscall.SIGPIPE:
				// Ignore random broken pipe
				traefiklog.FromContext(ctx).Infof("signal %q intercepted", sig.String())
			}
		case err = <-app.Errors():
			traefiklog.FromContext(ctx).WithError(err).Error("app error")
			break signalLoop
		case <-ctx.Done():
			break signalLoop
		}
	}

	go func() {
		signal.Stop(signals)
		close(signals)
	}()

	stopErr := app.Stop(ctx)
	if err != nil {
		return err
	}

	return stopErr
}

func (app *App) Stop(ctx context.Context) error {
	traefiklog.FromContext(ctx).Infof("app gracefully shutting down...")

	go func() {
		for range app.errors {
			// drain errors
		}
	}()

	// 1. interrupt daemons and wait for all daemons to complete
	app.cancel()
	app.daemonWg.Wait()

	// 2. stop grpc and http server
	defer app.serverWg.Wait()
	gr := &multierror.Group{}
	if app.http != nil {
		gr.Go(func() error { return tcp.Shutdown(ctx, app.http) })
	}

	if app.grpc != nil {
		gr.Go(func() error { return tcp.Shutdown(ctx, app.grpc) })
	}

	err := gr.Wait().ErrorOrNil()
	if err != nil {
		traefiklog.FromContext(ctx).WithError(err).Errorf("app could not shut down gracefully")
		return err // something went wrong while shutting down
	}

	traefiklog.FromContext(ctx).Infof("app gracefully shutted down")
	return nil // completed normally
}

func (app *App) Close() (err error) {
	app.closeOnce.Do(func() {
		close(app.errors)
		gr := &multierror.Group{}
		if app.http != nil {
			gr.Go(func() error { return tcp.Close(app.http) })
		}

		if app.grpc != nil {
			gr.Go(func() error { return tcp.Close(app.grpc) })
		}

		if app.watcher != nil {
			gr.Go(app.watcher.Close)
		}

		for _, daemon := range app.daemons {
			daemon := daemon
			gr.Go(daemon.Close)
		}

		err = gr.Wait().ErrorOrNil()
	})
	return
}

func (app *App) Errors() <-chan error {
	return app.errors
}

func (app *App) IsReady() bool {
	if !app.isReady {
		return false
	}

	gr := &multierror.Group{}
	for _, chk := range app.readinessCheck {
		gr.Go(chk.Check)
	}

	return gr.Wait().ErrorOrNil() == nil
}
