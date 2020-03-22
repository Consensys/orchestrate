package app

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/tcp"
)

type App struct {
	http *http.EntryPoints
	grpc *grpc.EntryPoint

	watcher configwatcher.Watcher

	cancel func()

	wg *sync.WaitGroup
}

func New(watcher configwatcher.Watcher, httpEps *http.EntryPoints, grpcEp *grpc.EntryPoint) *App {
	return &App{
		wg:      &sync.WaitGroup{},
		watcher: watcher,
		http:    httpEps,
		grpc:    grpcEp,
	}
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
