package configwatcher

import (
	"context"
	"time"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/aggregator"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/poll"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/use-cases"
)

const (
	InternalProviderName = "internal"
	ChainsProxyProvider  = "chains-proxy"
)

func NewProvider(
	getChains usecases.GetChains,
	cfg *dynamic.Configuration,
	refresh time.Duration,
) provider.Provider {
	prvdr := aggregator.New()
	prvdr.AddProvider(NewInternalProvider(cfg))
	prvdr.AddProvider(NewChainsProxyProvider(getChains, refresh))
	return prvdr
}

func NewInternalProvider(cfg *dynamic.Configuration) provider.Provider {
	return static.New(dynamic.NewMessage(InternalProviderName, cfg))
}

func NewChainsProxyProvider(getChains usecases.GetChains, refresh time.Duration) provider.Provider {
	poller := func(ctx context.Context) (provider.Message, error) {
		chains, err := getChains.Execute(ctx, "", nil)
		if err != nil {
			return nil, err
		}

		return dynamic.NewMessage(ChainsProxyProvider, newProxyConfig(chains)), nil
	}
	return poll.New(poller, refresh)
}
