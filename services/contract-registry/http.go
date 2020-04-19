package contractregistry

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
	svc "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	httpservice "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/service/http"
)

func NewHTTPBuilder(
	staticCfg *traefikstatic.Configuration,
	jwt, key auth.Checker,
	multitenancy bool,
	service svc.ContractRegistryServer,
	reg metrics.HTTP,
) (router.Builder, error) {
	builder := dynrouter.NewBuilder(staticCfg, nil)

	var err error
	// Create HTTP Handler Builder
	builder.Handler, err = NewHandlerBuilder(service)
	if err != nil {
		return nil, err
	}

	// Create Middleware Builder
	builder.Middleware, err = NewMiddlewareBuilder(jwt, key, multitenancy)
	if err != nil {
		return nil, err
	}

	builder.Metrics = metricsmid.NewBuilder(reg)

	return builder, nil
}

func NewHandlerBuilder(service svc.ContractRegistryServer) (handler.Builder, error) {
	builder := dynhandler.NewBuilder()

	// Add Builder for Contract-Registry API
	builder.AddBuilder(
		reflect.TypeOf(&dynamic.Contracts{}),
		httpservice.NewBuilder(service),
	)

	return builder, nil
}

func NewMiddlewareBuilder(jwt, key auth.Checker, multitenancy bool) (middleware.Builder, error) {
	builder := dynmid.NewBuilder()

	// Add Authentication Middleware
	builder.AddBuilder(
		reflect.TypeOf(&dynamic.Auth{}),
		authmid.NewBuilder(jwt, key, multitenancy),
	)

	return builder, nil
}
