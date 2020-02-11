package api

import (
	"context"
	"net/http"
	"sync"

	"github.com/containous/traefik/v2/pkg/api"
	"github.com/containous/traefik/v2/pkg/config/runtime"
	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/chains"
)

const component = "chain-registry.store.api"

var (
	initOnce = &sync.Once{}
)

// Initialize API handlers
func Init(ctx context.Context) {
	initOnce.Do(func() {
		chains.Init(ctx)
	})
}

type Builder func(config *runtime.Configuration) http.Handler

// NewBuilder returns a http.Handler builder based on runtime.Configuration
func NewBuilder(staticConfig *static.Configuration) Builder {
	return func(configuration *runtime.Configuration) http.Handler {
		router := mux.NewRouter()

		// Append Traefik API routes
		if staticConfig.API != nil {
			api.New(*staticConfig, configuration).Append(router)
		}

		// Append Chain-Registry routes
		chains.GlobalHandler().Append(router)

		return router
	}
}
