package api

import (
	"context"
	"net/http"
	"sync"

	"github.com/containous/traefik/v2/pkg/api"
	"github.com/containous/traefik/v2/pkg/config/runtime"
	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/gorilla/mux"
)

var (
	handler  *Handler
	initOnce = &sync.Once{}
)

// Initialize API handlers
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if handler != nil {
			return
		}

		// Set Chain-Registry handler
		handler = New()
	})
}

// NewBuilder returns a http.Handler builder based on runtime.Configuration
func NewBuilder(staticConfig *static.Configuration) Builder {
	return func(configuration *runtime.Configuration) http.Handler {
		router := mux.NewRouter()

		// Append Traefik API routes
		if staticConfig.API != nil {
			api.New(*staticConfig, configuration).Append(router)
		}

		// Append Chain-Registry routes
		handler.Append(router)

		return router
	}
}
