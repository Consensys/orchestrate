package app

import (
	"context"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/app/grpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/app/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/app/infra"
)

// App is main application object
type App struct {
	closeOnce *sync.Once
	done      chan struct{}
}

// New creates a new application
func New() *App {
	// We set a cancellable context so we can possibly abort application from within the application
	return &App{
		done:      make(chan struct{}),
		closeOnce: &sync.Once{},
	}
}

// Init application
func (app *App) init() {
	infra.Init()
	grpc.Init()
	http.Init()
}

// Run application
func (app *App) Run() {
	// Init
	app.init()

	// Create tcp listener
	l, err := net.Listen("tcp", viper.GetString("http.hostname"))
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"http.hostname": viper.GetString("http.hostname"),
		}).Fatalf("Can not listen tcp connection")
	}

	// Create MUX dispatcher to properly dispatch GRPC traffic and HTTP traffic
	tcpMux := cmux.New(l)
	grpcL := tcpMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := tcpMux.Match(cmux.HTTP1Fast())

	// Start serving
	wait := &sync.WaitGroup{}
	wait.Add(2)

	// Start GRPC server
	go func() {
		err := grpc.Server().Serve(grpcL)
		log.WithError(err).Info("app: grpc server stopped")
		wait.Done()
	}()

	// Start HTTP Server
	go func() {
		err := http.Server().Serve(httpL)
		log.WithError(err).Info("app: http server stopped")
		wait.Done()
	}()

	// Serve
	go func() {
		err := tcpMux.Serve()
		log.WithError(err).Fatal("app: tcpMux server stopped")
		wait.Done()
	}()

	log.WithFields(log.Fields{
		"http.hostname": l.Addr(),
	}).Infof("app: ready to receive connections")

	// Wait for server to properly close
	wait.Wait()

	// Close infra
	infra.Close()

	// We indicate that application has stopped running
	close(app.done)
	log.Infof("app: gracefully closed")
}

// Close gracefully stops the application
func (app *App) Close(ctx context.Context) {
	app.closeOnce.Do(func() {
		log.Warn("app: closing...")
		go grpc.Close(ctx)
		go http.Close(ctx)
	})
}

// Done return a channel indicating that application has gracefully closed
func (app *App) Done() <-chan struct{} {
	return app.done
}
