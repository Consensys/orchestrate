package api

import (
	"context"
	"net/http"
	"sync"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/swagger"

	"github.com/containous/traefik/v2/pkg/api"
	"github.com/containous/traefik/v2/pkg/config/runtime"
	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/chains"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/faucets"
)

//go:generate swag init --dir . --generalInfo exported.go --output ../../../public/swagger-specs/types/chain-registry
//go:generate rm ../../../public/swagger-specs/types/chain-registry/docs.go ../../../public/swagger-specs/types/chain-registry/swagger.yaml

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

const (
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
		faucets.Init(ctx)
		swagger.Init(swaggerSpecsPath, swaggerUIPath)
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
		swagger.GlobalHandler().Append(router)
		chains.GlobalHandler().Append(router)
		faucets.GlobalHandler().Append(router)

		return router
	}
}
