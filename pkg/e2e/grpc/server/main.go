package main

import (
	"context"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/grpc/server"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/healthcheck"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	defer log.Infof("Received: %v", in.Name)
	switch in.Name {
	case "":
		return &helloworld.HelloReply{}, errors.InvalidParameterError("no name provided").SetComponent("e2e.grpc.server")
	case "test-panic":
		panic("name made me panic")
	default:
		// Simulate io time
		time.Sleep(50 * time.Millisecond)
		return &helloworld.HelloReply{Message: "Hello " + in.Name}, nil
	}
}

var (
	app       *common.App
	startOnce = &sync.Once{}
)

func init() {
	// Create app
	app = common.NewApp()
}

// Run application
func Start(ctx context.Context) {
	startOnce.Do(func() {
		// Set log level to debug
		log.SetLevel(log.DebugLevel)

		// Initialize gRPC server
		grpcserver.AddEnhancers(
			func(s *grpc.Server) *grpc.Server {
				helloworld.RegisterGreeterServer(s, &server{})
				return s
			},
		)
		grpcserver.Init(ctx)

		// Initialize HTTP server
		http.Init(ctx)
		http.Enhance(healthcheck.HealthCheck(app))

		// Indicate that application is ready
		app.SetReady(true)

		// Start listening
		err := grpcserver.ListenAndServe()
		if err != nil {
			log.WithError(err).Error("main: error listening")
		}
	})
}

// Close gracefully stops the application
func Close(ctx context.Context) {
	log.Warn("app: closing...")
	_ = grpcserver.GracefulStop(ctx)
}

func main() {
	// Process signals
	sig := utils.NewSignalListener(func(signal os.Signal) { Close(context.Background()) })
	defer sig.Close()

	// Start application
	Start(context.Background())
}
