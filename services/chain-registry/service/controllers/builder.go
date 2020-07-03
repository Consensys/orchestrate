package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	chainsctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers/chains"
	faucetctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers/faucets"
)

//go:generate swag init --dir . --generalInfo builder.go --output ../../../../public/swagger-specs/services/chain-registry
//go:generate rm ../../../../public/swagger-specs/services/chain-registry/docs.go ../../../../public/swagger-specs/services/chain-registry/swagger.yaml

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

type Builder interface {
	Build(ctx context.Context, _ string, configuration interface{}, respModifier func(*http.Response) error) (http.Handler, error)
}

type builder struct {
	chainCtrl  chainsctrl.Controller
	faucetCtrl faucetctrl.Controller
}

func NewBuilder(chainCtrl chainsctrl.Controller, faucetCtrl faucetctrl.Controller) Builder {
	return &builder{
		chainCtrl:  chainCtrl,
		faucetCtrl: faucetCtrl,
	}
}

func (b *builder) Build(ctx context.Context, _ string, configuration interface{}, respModifier func(*http.Response) error) (http.Handler, error) {
	cfg, ok := configuration.(*dynamic.Chains)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type (expected %T but got %T)", cfg, configuration)
	}

	router := mux.NewRouter()
	b.chainCtrl.Append(router)
	b.faucetCtrl.Append(router)

	return router, nil
}
