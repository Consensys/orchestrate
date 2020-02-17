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

const (
	component        = "chain-registry.store.api"
	swaggerUIPath    = "./public/swagger-ui"
	swaggerSpecsPath = "./public/swagger-specs/types/chain-registry/swagger.json"
)

var (
	initOnce = &sync.Once{}
)

// Init Initialize API handlers
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

		// Append Swagger routes
		buildSwagger(router)

		return router
	}
}

// @title Chain Registry API
// @version 2.0
// @description PegaSys Orchestrate Chain Registry API. Enables dynamic management of chains

// @contact.name Contact PegaSys Orchestrate
// @contact.url https://pegasys.tech/contact/
// @contact.email support@pegasys.tech

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// @securityDefinitions.apikey JWTAuth
// @in header
// @name Authorization

func buildSwagger(router *mux.Router) {
	fs := http.FileServer(http.Dir(swaggerUIPath))
	router.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	router.HandleFunc("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, swaggerSpecsPath)
	})
}
