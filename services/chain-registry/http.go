package chainregistry

import (
	"reflect"

	traefikstatic "github.com/containous/traefik/v2/pkg/config/static"
	"github.com/dgraph-io/ristretto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler"
	dynhandler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/handler/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware"
	authmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/auth"
	dynmid "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/dynamic"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/middleware/ratelimit"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router"
	dynrouter "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/router/dynamic"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

func NewHTTPBuilder(
	staticCfg *traefikstatic.Configuration,
	jwt, key auth.Checker,
	multitenancy bool,
	s store.ChainRegistryStore,
	ec ethclient.ChainLedgerReader,
) (router.Builder, error) {
	builder := dynrouter.NewBuilder(staticCfg, nil)
	var err error
	// Create Service Builder
	builder.Handler, err = NewHandlerBuilder(staticCfg, s, ec)
	if err != nil {
		return nil, err
	}

	// Create Middleware Builder
	builder.Middleware, err = NewMiddlewareBuilder(jwt, key, multitenancy)
	if err != nil {
		return nil, err
	}

	return builder, nil
}

func NewHandlerBuilder(
	staticCfg *traefikstatic.Configuration,
	s store.ChainRegistryStore,
	ec ethclient.ChainLedgerReader,
) (handler.Builder, error) {
	builder := dynhandler.NewBuilder()

	// Chain-Registry API
	builder.AddBuilder(
		reflect.TypeOf(&dynamic.Chains{}),
		chainregistry.NewBuilder(s, ec),
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

func NewMiddlewareBuilder(jwt, key auth.Checker, multitenancy bool) (middleware.Builder, error) {
	builder := dynmid.NewBuilder()

	// RateLimit
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, err
	}
	builder.AddBuilder(
		reflect.TypeOf(&dynamic.RateLimit{}),
		ratelimit.NewBuilder(ratelimit.NewManager(cache)),
	)

	// Auth
	builder.AddBuilder(
		reflect.TypeOf(&dynamic.Auth{}),
		authmid.NewBuilder(jwt, key, multitenancy),
	)

	return builder, nil
}
