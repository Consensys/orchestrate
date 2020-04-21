package http

import (
	"reflect"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler"
	dynhandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/proxy"
	ctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers"
)

func newHandlerBuilder(
	builderCtrl ctrl.Builder,
	staticCfg *traefikstatic.Configuration,
) (handler.Builder, error) {
	builder := dynhandler.NewBuilder()

	// Chain-Registry API
	builder.AddBuilder(
		reflect.TypeOf(&dynamic.Chains{}),
		builderCtrl,
	)

	// ReverseProxy
	proxyBuilder, err := proxy.NewBuilder(staticCfg, nil)
	if err != nil {
		return nil, err
	}
	builder.AddBuilder(
		reflect.TypeOf(&dynamic.ReverseProxy{}),
		proxyBuilder,
	)

	return builder, nil
}
