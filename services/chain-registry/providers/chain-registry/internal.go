package internal

import (
	"math"

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
		ProviderName:  "internal",
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

	i.apiConfiguration(cfg)

	return cfg
}

func (i *Provider) apiConfiguration(cfg *dynamic.Configuration) {
	// Register api
	cfg.HTTP.Routers["api"] = &dynamic.Router{
		EntryPoints: []string{"orchestrate"},
		Service:     "api@internal",
		Priority:    math.MaxInt32 - 1,
		Rule:        "PathPrefix(`/{tenantID}`)",
		Middlewares: []string{"orchestrate-auth"},
	}
	cfg.HTTP.Services["api"] = &dynamic.Service{}
}
