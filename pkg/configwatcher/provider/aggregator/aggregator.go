package aggregator

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
)

// Provider aggregates providers.
type Provider struct {
	providers []provider.Provider
}

// New create new aggregator
func New() *Provider {
	return &Provider{}
}

// AddProvider adds a provider in the providers map.
func (p *Provider) AddProvider(prvdr provider.Provider) {
	p.providers = append(p.providers, prvdr)
}

// Provide calls the provide method of every providers
func (p *Provider) Provide(ctx context.Context, msgs chan<- provider.Message) error {
	wg := &sync.WaitGroup{}
	wg.Add(len(p.providers))
	for _, prvdr := range p.providers {
		go func(prvdr provider.Provider) {
			defer wg.Done()
			err := prvdr.Provide(ctx, msgs)
			if err != nil {
				log.FromContext(ctx).WithError(err).Errorf("Provider %T failed", prvdr)
			}
		}(prvdr)
	}
	wg.Wait()

	return nil
}
