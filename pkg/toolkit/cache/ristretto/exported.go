package ristretto

import (
	"context"
	"sync"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	ristretto2 "github.com/consensys/orchestrate/pkg/toolkit/cache/noncache"
	"github.com/dgraph-io/ristretto"

	"github.com/consensys/orchestrate/pkg/toolkit/cache"
	"github.com/spf13/viper"
)

const component = "service.cache-manager"

var (
	client   cache.Manager
	initOnce = &sync.Once{}
)

func Init(_ context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		logger := log.NewLogger().SetComponent(component)

		vipr := viper.GetViper()
		cfg := NewConfig(vipr)

		if cfg.CacheTTL == nil {
			logger.Info("disabled")
			client = ristretto2.NewNonCacheManager()
			return
		}

		// Set Client
		cCache, err := ristretto.NewCache(cfg.Cache)
		if err != nil {
			logger.Error("failed to initialize")
			client = ristretto2.NewNonCacheManager()
			return
		}
		client = NewCacheManager(component, cCache, cfg.CacheTTL)

		logger.WithField("ttl", cfg.CacheTTL).Info("ready")
	})
}

// GlobalClient returns global Client
func GlobalClient() cache.Manager {
	return client
}

// SetGlobalClient sets global Client
func SetGlobalClient(ec cache.Manager) {
	client = ec
}
