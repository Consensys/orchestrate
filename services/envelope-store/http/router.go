package http

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router"
	dynrouter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router/dynamic"
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/proto"
)

type RouterBuilder router.Builder

func NewRouterBuilder(
	srv svc.EnvelopeStoreServer,
	staticCfg Config,
	jwt, key auth.Checker,
	multitenancy bool,
) (RouterBuilder, error) {
	builder := dynrouter.NewBuilder(staticCfg, nil)

	var err error
	// Create HTTP Handler Builder
	builder.Handler, err = newHandlerBuilder(srv)
	if err != nil {
		return nil, err
	}

	// Create Middleware Builder
	builder.Middleware, err = newMiddlewareBuilder(jwt, key, multitenancy)
	if err != nil {
		return nil, err
	}

	return builder, nil
}


