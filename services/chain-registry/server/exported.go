package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/containous/alice"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/provider/acme"
	"github.com/containous/traefik/v2/pkg/provider/aggregator"
	traefiktls "github.com/containous/traefik/v2/pkg/tls"
	"github.com/dgraph-io/ristretto"
	"github.com/spf13/viper"
	authjwt "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/jwt"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/key"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api"
	authmw "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/middlewares/auth"
	ratelimitermw "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/middlewares/ratelimiter"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/providers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/server/router"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
)

var (
	staticConfig *static.Configuration
	svr          *Server
	initOnce     = &sync.Once{}
)

func initACMEProvider(c *static.Configuration, providerAggregator *aggregator.ProviderAggregator, tlsManager *traefiktls.Manager) []*acme.Provider {
	challengeStore := acme.NewLocalChallengeStore()
	localStores := map[string]*acme.LocalStore{}

	var resolvers []*acme.Provider
	for name, resolver := range c.CertificatesResolvers {
		if resolver.ACME == nil {
			continue
		}

		if localStores[resolver.ACME.Storage] == nil {
			localStores[resolver.ACME.Storage] = acme.NewLocalStore(resolver.ACME.Storage)
		}

		p := &acme.Provider{
			Configuration:  resolver.ACME,
			Store:          localStores[resolver.ACME.Storage],
			ChallengeStore: challengeStore,
			ResolverName:   name,
		}

		if err := providerAggregator.AddProvider(p); err != nil {
			log.WithoutContext().Errorf("Unable to add ACME provider to the providers list: %v", err)
			continue
		}
		p.SetTLSManager(tlsManager)
		if p.TLSChallenge != nil {
			tlsManager.TLSAlpnGetter = p.GetTLSALPNCertificate
		}
		p.SetConfigListenerChan(make(chan dynamic.Configuration))
		resolvers = append(resolvers, p)
	}
	return resolvers
}

func initOrchestrateMiddlewares(ctx context.Context) map[string]func(string) alice.Constructor {
	authkey.Init(ctx)
	authjwt.Init(ctx)

	middlewares := make(map[string]func(string) alice.Constructor)

	// Set Authentication middleware
	middlewares["orchestrate-auth"] = func(routerName string) alice.Constructor {
		return func(next http.Handler) (http.Handler, error) {
			return authmw.New(
				authjwt.GlobalAuth(),
				authkey.GlobalAuth(),
				viper.GetBool(multitenancy.EnabledViperKey),
				next), nil
		}
	}

	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	mngr := ratelimitermw.NewManager(cache)
	middlewares["orchestrate-ratelimit"] = func(routerName string) alice.Constructor {
		return func(next http.Handler) (http.Handler, error) {
			return ratelimitermw.New(
				mngr,
				time.Second,
				30*time.Second,
				5*time.Second,
				routerName,
				next), nil
		}
	}

	return middlewares
}

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if svr != nil {
			return
		}

		// Initialize API
		api.Init(ctx)

		// Initialize providers
		providers.Init(ctx)

		tlsManager := traefiktls.NewManager()

		acmeProviders := initACMEProvider(staticConfig, providers.ProviderAggregator(), tlsManager)

		serverEntryPointsTCP := make(TCPEntryPoints)
		var err error
		for entryPointName, config := range staticConfig.EntryPoints {
			logCtx := log.With(context.Background(), log.Str(log.EntryPointName, entryPointName))
			serverEntryPointsTCP[entryPointName], err = NewTCPEntryPoint(logCtx, config)
			if err != nil {
				log.WithoutContext().WithError(err).Fatalf("error while building entryPoint %s", entryPointName)
			}
			serverEntryPointsTCP[entryPointName].RouteAppenderFactory = router.NewRouteAppenderFactory(staticConfig, entryPointName, nil)
		}

		// Initialize custom orchestrate middleware
		middlewares := initOrchestrateMiddlewares(ctx)

		svr = NewServer(staticConfig, providers.ProviderAggregator(), serverEntryPointsTCP, tlsManager, api.NewBuilder(staticConfig), middlewares)

		resolverNames := map[string]struct{}{}
		for _, p := range acmeProviders {
			resolverNames[p.ResolverName] = struct{}{}
			svr.AddListener(p.ListenConfiguration)
		}
	})
}

// SetGlobalStaticConfig set Traefik static configuration for global server
func SetGlobalStaticConfig(config *static.Configuration) {
	staticConfig = config
}

// SetGlobalServer set global server
func SetGlobalServer(server *Server) {
	svr = server
}

// GlobalServer returns global server
func GlobalServer() *Server {
	return svr
}
