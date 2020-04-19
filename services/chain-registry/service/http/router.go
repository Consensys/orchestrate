package http

import (
	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	metricsmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router"
	dynrouter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics"
	ctrl "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/service/controllers"
)

func NewHTTPRouterBuilder(
	builderCtrl ctrl.Builder,
	staticCfg *traefikstatic.Configuration,
	jwt, key auth.Checker,
	multitenancy bool,
	reg metrics.HTTP,
) (router.Builder, error) {
	builder := dynrouter.NewBuilder(staticCfg, nil)
	var err error
	// Create Service Controller
	builder.Handler, err = newHandlerBuilder(builderCtrl, staticCfg)
	if err != nil {
		return nil, err
	}

	// Create Middleware Controller
	builder.Middleware, err = newMiddlewareBuilder(jwt, key, multitenancy)
	if err != nil {
		return nil, err
	}

	builder.Metrics = metricsmid.NewBuilder(reg)

	return builder, nil
}
