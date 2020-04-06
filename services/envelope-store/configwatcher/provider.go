package configwatcher

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/aggregator"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/configwatcher/provider/static"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/config/dynamic"
)

const (
	InternalProviderName = "internal"
)

func NewProvider(
	cfg *dynamic.Configuration,
) provider.Provider {
	prvdr := aggregator.New()
	prvdr.AddProvider(NewInternalProvider(cfg))
	return prvdr
}

func NewInternalProvider(cfg *dynamic.Configuration) provider.Provider {
	return static.New(dynamic.NewMessage(InternalProviderName, cfg))
}
