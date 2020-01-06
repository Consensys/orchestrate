package nodes

import (
	"context"
	"sync"

	"github.com/containous/traefik/v2/pkg/provider"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

var (
	provide  provider.Provider
	initOnce = &sync.Once{}
)

// Init initializes a ChainRegistry store
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if provide != nil {
			return
		}

		store.Init(ctx)
		registry := store.GlobalStoreRegistry()
		provide = NewProvider(viper.GetString(store.TypeViperKey), registry, viper.GetDuration(ProviderRefreshIntervalViperKey))

	})
}

// SetGlobalRegistry sets global a chain-registry store
func SetGlobalProvider(p provider.Provider) {
	provide = p
}

// GlobalRegistry returns global a chain-registry store
func GlobalProvider() provider.Provider {
	return provide
}
