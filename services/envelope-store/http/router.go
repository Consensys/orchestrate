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

	// Create HTTP Handler Builder
	builder.Handler = newHandlerBuilder(srv)

	// Create Middleware Builder
	builder.Middleware = newMiddlewareBuilder(jwt, key, multitenancy)

	return builder, nil
}
