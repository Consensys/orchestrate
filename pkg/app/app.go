package app

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tcp"
)

type App struct {
	cfg *Config

	http *http.EntryPoints
	grpc *grpc.EntryPoint

	watcher configwatcher.Watcher

	cancel func()

	wg *sync.WaitGroup

	metrics metrics.Registry
}

func New(
	cfg *Config,
	prvdr provider.Provider,
	httpBuilder router.Builder,
	grpcBuilder grpcserver.Builder,
	reg metrics.Registry,
) (*App, error) {
	// Create HTTP EntryPoints
	httpEps := http.NewEntryPoints(cfg.HTTP.EntryPoints, httpBuilder, reg.TCP())

	// Create watcher
	watcher := configwatcher.New(
		cfg.Watcher,
		prvdr,
		dynamic.Merge,
		// Listeners
		[]func(context.Context, interface{}) error{
			httpEps.Switch,
			func(_ context.Context, cfg interface{}) error {
				if dynCfg, ok := cfg.(*dynamic.Configuration); ok {
					return reg.HTTP().Switch(dynCfg)
				}
				return nil
			},
		},
	)

	var grpcEp *grpc.EntryPoint
	if grpcBuilder != nil && cfg.GRPC != nil {
		grpcServer, err := grpcBuilder.Build(context.Background(), "", cfg.GRPC.Static)
		if err != nil {
			return nil, err
		}

		grpcEp = grpc.NewEntryPoint("", cfg.GRPC.EntryPoint, grpcServer, reg.TCP())
	}

	return &App{
		cfg:     cfg,
		wg:      &sync.WaitGroup{},
		watcher: watcher,
		http:    httpEps,
		grpc:    grpcEp,
		metrics: reg,
	}, nil
}

func (app *App) Start(ctx context.Context) error {
	cancelableCtx, cancel := context.WithCancel(context.Background())
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

	return nil
}

func (app *App) Stop(ctx context.Context) error {
	log.FromContext(ctx).Infof("gracefully shutting down application")
	app.cancel()

	app.wg.Add(2)
	go func() {
		if app.http != nil {
			_ = tcp.Shutdown(ctx, app.http)
			_ = tcp.Close(app.http)
		}
		app.wg.Done()
	}()

	go func() {
		if app.grpc != nil {
			_ = tcp.Shutdown(ctx, app.grpc)
			_ = tcp.Close(app.grpc)
		}
		app.wg.Done()
	}()

	closed := make(chan struct{})
	go func() {
		app.wg.Wait()
		close(closed)
	}()

	select {
	case <-closed:
		return nil // completed normally
	case <-ctx.Done():
		return ctx.Err() // timed out
	}
}
