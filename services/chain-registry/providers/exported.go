package providers

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/provider/aggregator"
	internal "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/providers/chain-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/providers/test"
)

var (
	providerAggregator *aggregator.ProviderAggregator
	initOnce           = &sync.Once{}
)

// Initialize provider aggregator
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if providerAggregator != nil {
			return
		}

		providerAggregator = &aggregator.ProviderAggregator{}

		err := providerAggregator.AddProvider(internal.New())
		if err != nil {
			log.WithoutContext().WithError(err).Fatalf("error adding internal provider")
		}

		// TODO: test provider should be replaced
		err = providerAggregator.AddProvider(test.New())
		if err != nil {
			log.WithoutContext().WithError(err).Fatalf("error adding test provider")
		}
	})
}

// Return global provider aggregator
func ProviderAggregator() *aggregator.ProviderAggregator {
	return providerAggregator
}

// Set global provider aggregator
func SetGlobalProviderAggregator(provider *aggregator.ProviderAggregator) {
	providerAggregator = provider
}
