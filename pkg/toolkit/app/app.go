package app

import (
	"context"
	"encoding/json"
	nethttp "net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/config/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/provider"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/configwatcher/provider/aggregator"
	httphandler "github.com/consensys/orchestrate/pkg/toolkit/app/http/handler/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/handler/healthcheck"
	httpmetrics "github.com/consensys/orchestrate/pkg/toolkit/app/http/metrics"
	httpmid "github.com/consensys/orchestrate/pkg/toolkit/app/http/middleware/dynamic"
	metricsmid "github.com/consensys/orchestrate/pkg/toolkit/app/http/middleware/metrics"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/router"
	httprouter "github.com/consensys/orchestrate/pkg/toolkit/app/http/router/dynamic"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/app/metrics"
	metricsregistry "github.com/consensys/orchestrate/pkg/toolkit/app/metrics/registry"
	"github.com/consensys/orchestrate/pkg/toolkit/tcp"
	tcpmetrics "github.com/consensys/orchestrate/pkg/toolkit/tcp/metrics"
	"github.com/hashicorp/go-multierror"
	healthz "github.com/heptiolabs/healthcheck"
	"github.com/sirupsen/logrus"
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

	provider *aggregator.Provider
	watcher  configwatcher.Watcher

	metricReg metrics.Registry

	daemons        []Daemon
	readinessCheck []*healthcheck.Checker

	logger *log.Logger

	cancel func()

	daemonWg sync.WaitGroup
	serverWg sync.WaitGroup

	errors chan error

	closeOnce sync.Once

	isReady bool
}

func New(cfg *Config, opts ...Option) (*App, error) {
	err := log.ConfigureLogger(cfg.Log, logrus.StandardLogger())
	if err != nil {
		return nil, err
	}

	reg := metricsregistry.NewMetricRegistry()

	httpBuilder := httprouter.NewBuilder(cfg.HTTP.TraefikStatic())
	httpBuilder.Handler = httphandler.NewBuilder()
	httpBuilder.Middleware = httpmid.NewBuilder()

	var httpMidMetrics httpmetrics.HTTPMetrics
	if cfg.Metrics.IsActive(httpmetrics.ModuleName) {
		httpMidMetrics = httpmetrics.NewHTTPMetrics(nil)
	} else {
		httpMidMetrics = httpmetrics.NewHTTPNopMetrics(nil)
	}

	reg.Add(httpMidMetrics)
	httpBuilder.Metrics = metricsmid.NewBuilder(httpMidMetrics)

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
		watcher,
		reg,
		log.NewLogger().SetComponent("application"),
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
	watcher configwatcher.Watcher,
	reg metrics.Registry,
	logger *log.Logger,
) *App {
	return &App{
		cfg:         cfg,
		httpBuilder: httpBuilder,
		watcher:     watcher,
		logger:      logger,
		metricReg:   reg,
		errors:      make(chan error),
	}
}

func (app *App) init(_ context.Context) error {
	conf, err := json.Marshal(app.cfg)
	if err != nil {
		return err
	}
	app.logger.WithField("conf", string(conf)).Debug("loaded app configuration")
	app.logger.WithField("metrics", app.cfg.Metrics.Modules()).Info("activated metric modules")
	var tcpreg tcpmetrics.TPCMetrics
	if app.cfg.HTTP != nil {
		if app.cfg.Metrics.IsActive(tcpmetrics.ModuleName) {
			tcpreg = tcpmetrics.NewTCPMetrics(nil)
		} else {
			tcpreg = tcpmetrics.NewTCPNopMetrics(nil)
		}

		app.metricReg.Add(tcpreg)
	}

	// Create HTTP EntryPoints
	if app.cfg.HTTP != nil {
		app.http = http.NewEntryPoints(app.cfg.HTTP.EntryPoints, app.httpBuilder, tcpreg)

		// Add Listeners for HTTP
		app.watcher.AddListener(app.http.Switch)
		app.watcher.AddListener(
			func(_ context.Context, cfg interface{}) error {
				if dynCfg, ok := cfg.(*dynamic.Configuration); ok {
					if err := app.metricReg.SwitchDynConfig(dynCfg); err != nil {
						return err
					}
				}
				return nil
			},
		)
	}

	return nil
}

func (app *App) HTTP() router.Builder {
	return app.httpBuilder
}

func (app *App) MetricRegistry() metrics.Registry {
	return app.metricReg
}

func (app *App) Logger() *log.Logger {
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

	app.logger.Debug("starting...")

	if app.http != nil {
		app.serverWg.Add(1)
		go func() {
			for err := range app.http.ListenAndServe(ctx) {
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
			app.errors <- err
			app.daemonWg.Done()
		}()
	}

	app.daemonWg.Add(len(app.daemons))
	for _, daemon := range app.daemons {
		go func(daemon Daemon) {
			err := daemon.Run(cancelableCtx)
			app.errors <- err
			app.daemonWg.Done()
		}(daemon)
	}

	app.isReady = true

	app.logger.Info("started")
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
				app.logger.WithField("sig", sig.String()).Debug("signal intercepted")
				break signalLoop
			}
		case err = <-app.Errors():
			if err != nil && err != context.Canceled && err != nethttp.ErrServerClosed {
				app.logger.WithError(err).Error("exits with errors")
			} else {
				app.logger.WithError(err).Info("exits")
			}
			break signalLoop
		case <-ctx.Done():
			app.logger.WithError(ctx.Err()).Info("exits gracefully")
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
	app.logger.Debug("shutting down...")

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

	err := gr.Wait().ErrorOrNil()
	if err != nil {
		app.logger.WithError(err).Error("could not shut down gracefully")
		return err // something went wrong while shutting down
	}

	app.isReady = false
	app.logger.WithError(err).Info("gracefully shutted down")
	return nil // completed normally
}

func (app *App) Close() (err error) {
	app.closeOnce.Do(func() {
		close(app.errors)
		gr := &multierror.Group{}
		if app.http != nil {
			gr.Go(func() error { return tcp.Close(app.http) })
		}

		if app.watcher != nil {
			gr.Go(app.watcher.Close)
		}

		for _, daemon := range app.daemons {
			daemon := daemon
			gr.Go(daemon.Close)
		}

		app.isReady = false
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
