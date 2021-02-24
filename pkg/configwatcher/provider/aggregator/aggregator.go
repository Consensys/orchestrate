package aggregator

import (
	"context"
	"fmt"
	"sync"

	"github.com/ConsenSys/orchestrate/pkg/configwatcher/provider"
	"github.com/ConsenSys/orchestrate/pkg/log"
)

// Provider aggregates providers.
type Provider struct {
	providers []provider.Provider
	logger    *log.Logger
}

// New create new aggregator
func New() *Provider {
	return &Provider{
		logger: log.NewLogger().SetComponent("configwatcher"),
	}
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
				p.logger.WithField("provider", fmt.Sprintf("%T", prvdr)).
					WithError(err).Error("provider failed")
			}
		}(prvdr)
	}
	wg.Wait()

	return nil
}
