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
