package grpcserver

import (
	"context"
	"net"
	"net/http"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"google.golang.org/grpc"
)

const netClosingErrorMsg = "use of closed network connection"

// CMuxServer is a server that enables to simultaneously serve gRPC and HTTP traffic
type CMuxServer struct {
	grpc *grpc.Server
	http *http.Server
}

// NewCMuxServer creates a new CMuxServer
func NewCMuxServer(grpcsrv *grpc.Server, httpsev *http.Server) *CMuxServer {
	return &CMuxServer{
		grpc: grpcsrv,
		http: httpsev,
	}
}

// Listen announces on the local network address and accepts incoming connection
//
// The network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
func (s *CMuxServer) ListenAndServe(network, address string) error {
	// Open tcp connection
	lis, err := net.Listen(network, address)
	if err != nil {
		return errors.ConnectionError(err.Error()).SetComponent(component)
	}
	return s.Serve(lis)
}

// Serve accepts incoming connections on the listener lis
func (s *CMuxServer) Serve(lis net.Listener) error {
	// Create a multiplexer to dispatch traffic between gRPC and HTTP
	tcpMux := cmux.New(lis)
	grpcL := tcpMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := tcpMux.Match(cmux.HTTP1Fast())

	wait := sync.WaitGroup{}
	wait.Add(2)
	defer wait.Wait()

	// Start gRPC server
	go func() {
		if err := s.grpc.Serve(grpcL); err != nil && err != grpc.ErrServerStopped && !strings.Contains(err.Error(), netClosingErrorMsg) {
			e := errors.GRPCConnectionError(err.Error()).SetComponent(component)
			log.WithError(e).Warn("cmuxserver: error while serving GRPC")
		}
		wait.Done()
	}()

	// Start HTTP Server
	go func() {
		if err := s.http.Serve(httpL); err != nil && err != http.ErrServerClosed && !strings.Contains(err.Error(), netClosingErrorMsg) {
			e := errors.HTTPConnectionError(err.Error()).SetComponent(component)
			log.WithError(e).Warn("cmuxserver: error while serving HTTP")
		}
		wait.Done()
	}()

	// Serve
	if err := tcpMux.Serve(); err != nil && !strings.Contains(err.Error(), netClosingErrorMsg) {
		return errors.ConnectionError(err.Error()).SetComponent(component)
	}

	return nil
}

// Shutdown gracefully shuts down the server without interrupting any
// active connections.
func (s *CMuxServer) Shutdown(ctx context.Context) error {
	var e error
	common.InParallel(
		// GracefulStop gRPC server
		func() { s.grpc.GracefulStop() },
		// Stop HTTP server
		func() {
			err := s.http.Shutdown(ctx)
			if err != nil {
				e = errors.ConnectionError(err.Error()).SetComponent(component)
			}
		},
	)

	return e
}
