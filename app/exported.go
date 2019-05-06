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

var (
	app       *App
	startOnce = &sync.Once{}
)

func init() {
	// Create app
	app = NewApp()
}

// Init application
func initComponents() {
	infra.Init()
	grpc.Init()
	http.Init()
}

// Run application
func Start(ctx context.Context) {
	startOnce.Do(func() {
		// Init
		initComponents()

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
		wait.Add(3)

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
			log.WithError(err).Info("app: tcpMux server stopped")
			wait.Done()
		}()

		log.WithFields(log.Fields{
			"http.hostname": l.Addr(),
		}).Infof("app: ready to receive connections")

		// Indicate that application is ready
		// TODO: we need to update so ready can append when Consume has finished to Setup
		app.ready.Store(true)

		// Wait for server to properly close
		wait.Wait()

		// Close infra
		infra.Close()

		log.Infof("app: gracefully closed")
	})
}

// Close gracefully stops the application
func Close(ctx context.Context) {
	log.Warn("app: closing...")
	go grpc.Close(ctx)
	go http.Close(ctx)
}
