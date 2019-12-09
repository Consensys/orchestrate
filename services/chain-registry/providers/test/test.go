package test

import (
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/provider"
	"github.com/containous/traefik/v2/pkg/safe"
	"github.com/containous/traefik/v2/pkg/tls"
)

var _ provider.Provider = (*Provider)(nil)

// Provider is a provider.Provider implementation that provides the internal routers.
type Provider struct{}

// New creates a new instance of the internal provider.
func New() *Provider {
	return &Provider{}
}

// Provide allows the provider to provide configurations to traefik using the given configuration channel.
func (i *Provider) Provide(configurationChan chan<- dynamic.Message, _ *safe.Pool) error {
	configurationChan <- dynamic.Message{
		ProviderName:  "test",
		Configuration: i.createConfiguration(),
	}

	return nil
}

// Init the provider.
func (i *Provider) Init() error {
	return nil
}

func (i *Provider) createConfiguration() *dynamic.Configuration {
	cfg := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:     make(map[string]*dynamic.Router),
			Middlewares: make(map[string]*dynamic.Middleware),
			Services:    make(map[string]*dynamic.Service),
		},
		TCP: &dynamic.TCPConfiguration{
			Routers:  make(map[string]*dynamic.TCPRouter),
			Services: make(map[string]*dynamic.TCPService),
		},
		TLS: &dynamic.TLSConfiguration{
			Stores:  make(map[string]tls.Store),
			Options: make(map[string]tls.Options),
		},
	}

	i.testConfiguration(cfg)

	return cfg
}

func (i *Provider) testConfiguration(cfg *dynamic.Configuration) {
	cfg.HTTP.Services["test-node"] = &dynamic.Service{
		LoadBalancer: &dynamic.ServersLoadBalancer{
			Servers: []dynamic.Server{
				dynamic.Server{
					URL: "http://localhost:8545",
				},
			},
		},
	}

	cfg.HTTP.Services["infura-mainnet"] = &dynamic.Service{
		LoadBalancer: &dynamic.ServersLoadBalancer{
			Servers: []dynamic.Server{
				dynamic.Server{
					Scheme: "https",
					URL:    "https://mainnet.infura.io/v3/bfc9d6e51fbc4d3db54bea58d1094f9c",
				},
			},
		},
	}

	cfg.HTTP.Routers["test-node"] = &dynamic.Router{
		EntryPoints: []string{"http"},
		Service:     "test-node",
		// Priority:    math.MaxInt64 - 1,
		Rule: "Path(`/test-node`)",
	}

	cfg.HTTP.Routers["infura-mainnet"] = &dynamic.Router{
		EntryPoints: []string{"http"},
		// Middlewares: []string{"redirect-https"},
		Service: "infura-mainnet",
		// Priority:    math.MaxInt64 - 1,
		Rule: "Path(`/infura-mainnet`)",
	}
}
