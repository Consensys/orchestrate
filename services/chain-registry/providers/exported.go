package providers

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/provider/aggregator"
	internal "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/providers/chain-registry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/providers/nodes"
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

		nodes.Init()
		err = providerAggregator.AddProvider(nodes.GlobalProvider())
		if err != nil {
			log.WithoutContext().WithError(err).Fatalf("error adding registry provider")
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
