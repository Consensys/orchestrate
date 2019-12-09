package server

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/config/static"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/provider/acme"
	"github.com/containous/traefik/v2/pkg/provider/aggregator"
	"github.com/containous/traefik/v2/pkg/server/router"
	traefiktls "github.com/containous/traefik/v2/pkg/tls"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/providers"
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
			ctx := log.With(context.Background(), log.Str(log.EntryPointName, entryPointName))
			serverEntryPointsTCP[entryPointName], err = NewTCPEntryPoint(ctx, config)
			if err != nil {
				log.WithoutContext().WithError(err).Fatalf("error while building entryPoint %s", entryPointName)
			}
			serverEntryPointsTCP[entryPointName].RouteAppenderFactory = router.NewRouteAppenderFactory(*staticConfig, entryPointName, nil)
		}

		svr = NewServer(staticConfig, providers.ProviderAggregator(), serverEntryPointsTCP, tlsManager, api.NewBuilder(staticConfig))

		resolverNames := map[string]struct{}{}
		for _, p := range acmeProviders {
			resolverNames[p.ResolverName] = struct{}{}
			svr.AddListener(p.ListenConfiguration)
		}
	})
}

// SetGlobalStaticConfig set traefil static configuration for global server
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
