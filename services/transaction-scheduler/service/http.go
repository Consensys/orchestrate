package service

import (
	"reflect"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler"
	dynhandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware"
	authmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/auth"
	dynmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/dynamic"
	metricsmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router"
	dynrouter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/metrics"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/controllers"
)

func NewHTTPBuilder(
	staticCfg *traefikstatic.Configuration,
	jwt, key auth.Checker,
	multitenancy bool,
	ctrls *controllers.Builder,
	reg metrics.HTTP,
) (router.Builder, error) {
	builder := dynrouter.NewBuilder(staticCfg, nil)

	// Create Service Builder
	builder.Handler = newHandlerBuilder(ctrls)
	// Create Middleware Builder
	builder.Middleware = newMiddlewareBuilder(jwt, key, multitenancy)

	builder.Metrics = metricsmid.NewBuilder(reg)

	return builder, nil
}

func newHandlerBuilder(ctrls *controllers.Builder) handler.Builder {
	builder := dynhandler.NewBuilder()

	// Transaction API
	builder.AddBuilder(
		reflect.TypeOf(&dynamic.Transactions{}),
		ctrls,
	)

	return builder
}

func newMiddlewareBuilder(jwt, key auth.Checker, multitenancy bool) middleware.Builder {
	builder := dynmid.NewBuilder()

	// Auth
	builder.AddBuilder(
		reflect.TypeOf(&dynamic.Auth{}),
		authmid.NewBuilder(jwt, key, multitenancy),
	)

	return builder
}
