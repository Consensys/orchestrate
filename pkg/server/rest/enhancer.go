package rest

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Enhancer are functions that enhance net/http Multiplexers
type Enhancer func(ctx context.Context, mux *http.ServeMux, gwMux *runtime.ServeMux, conn *grpc.ClientConn) error

// ApplyEnhancers apply enhancers on a server
func ApplyEnhancers(mux *http.ServeMux, gwMux *runtime.ServeMux, conn *grpc.ClientConn, enhancers ...Enhancer) {
	// Enhance server
	for _, enhancer := range enhancers {
		err := enhancer(context.Background(), mux, gwMux, conn)
		if err != nil {
			log.WithError(err).Fatal("cannot start rest server")
		}
	}
}
