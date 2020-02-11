package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	grpcServer "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/grpc"
	"google.golang.org/grpc"
)

const (
	component        = "rest.server"
	swaggerUIPath    = "./public/swagger-ui"
	SwaggerSpecsPath = "./public/swagger-specs"
)

var (
	initOnce  = &sync.Once{}
	enhancers []Enhancer
	server    *http.Server
)

// Init initialize global gRPC server
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if server == nil {
			grpcConn, err := grpc.Dial(
				fmt.Sprintf("%s:%d", viper.GetString(grpcServer.HostnameViperKey), viper.GetUint(grpcServer.PortViperKey)),
				grpc.WithInsecure())
			if err != nil {
				log.WithError(err).Fatal("cannot start api server")
			}

			mux := http.NewServeMux()
			fs := http.FileServer(http.Dir(swaggerUIPath))
			mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

			gw := runtime.NewServeMux()
			ApplyEnhancers(mux, gw, grpcConn, enhancers...)
			mux.Handle("/", gw)

			server = &http.Server{
				Addr:         URL(),
				Handler:      mux,
				WriteTimeout: 10 * time.Second,
				ReadTimeout:  10 * time.Second,
			}
		}
	})
}

// AddEnhancers adds gRPC server enhancers that will be called at Init time
// Note that it should be called before Init()
func AddEnhancers(fns ...Enhancer) {
	enhancers = append(enhancers, fns...)
}

// GlobalServer return global gRPC server
func GlobalServer() *http.Server {
	return server
}

// SetGlobalServer sets global gRPC server
func SetGlobalServer(s *http.Server) {
	server = s
}

// ListenAndServe starts global server
func ListenAndServe() {
	if server == nil {
		log.Fatalf("%s: server is not initialized", component)
	}

	lis, err := net.Listen("tcp", URL())
	if err != nil {
		log.WithError(errors.ConnectionError(err.Error()).ExtendComponent(component)).
			WithFields(log.Fields{"grpc.url": URL()}).
			Error("failed to listen")
		return
	}
	log.Infof("%s: start serving on %q", URL(), component)

	// Serve requests
	err = server.Serve(lis)
	if err != nil && err != http.ErrServerClosed {
		log.WithError(errors.FromError(err).ExtendComponent(component)).
			WithFields(log.Fields{"rest.url": URL()}).
			Errorf("%s: failed to run REST server on URL %s", component, URL())
	} else {
		log.Infof("%s: server stopped", component)
	}
}

func StopServer(ctx context.Context) {
	err := server.Shutdown(ctx)
	if err != nil {
		log.WithError(err).Errorf("rest server failed to gracefully shutdown")
	}
}
