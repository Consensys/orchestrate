package metrics

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/julien-marchand/healthcheck"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

const component = "metrics"

var (
	initOnce = &sync.Once{}
	mux      = http.NewServeMux()
	server   *http.Server
)

// Init initialize global HTTP server
func Init(_ context.Context) {
	initOnce.Do(func() {
		if server != nil {
			return
		}

		// Initialize server
		server = &http.Server{}
		server.Addr = URL()
		server.Handler = mux
	})
}

// Enhance allows to register new handlers on Global Server ServeMux
func Enhance(enhancer ServeMuxEnhancer) {
	enhancer(mux)
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

	err = server.Serve(lis)
	if err != nil && err != http.ErrServerClosed {
		ierr := errors.HTTPConnectionError(err.Error()).SetComponent(component)
		log.WithError(ierr).Errorf("%s: error while listening", component)
	} else {
		log.Infof("%s: server gracefully stopped", component)
	}
}

func StartServer(ctx context.Context, cancel context.CancelFunc, isAlive, isReady healthcheck.Check) {
	// Initialize server
	Init(ctx)

	// Register Healthcheck
	Enhance(Enhancer(isAlive, isReady))

	// Start Listening
	ListenAndServe()

	cancel()
}

func StopServer(ctx context.Context) {
	err := server.Shutdown(ctx)
	if err != nil {
		log.WithError(err).Errorf("metrics server failed to gracefully shutdown")
	}
}

// GlobalServer return global HTTP server
func GlobalServer() *http.Server {
	return server
}

// SetGlobalServer sets global HTTP server
func SetGlobalServer(s *http.Server) {
	server = s
}
