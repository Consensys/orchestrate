package quorumkeymanager

import (
	"sync"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	qkm "github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/spf13/viper"
)

const component = "quorum-key-manager.client"

var (
	client      qkm.KeyManagerClient
	storeNameID string
	initOnce    = &sync.Once{}
)

func Init() {
	initOnce.Do(func() {
		if client != nil {
			return
		}
		vipr := viper.GetViper()

		logger := log.NewLogger().SetComponent(component)
		cfg := NewConfigFromViper(vipr)
		if cfg.URL != "" {
			httpClient, err := NewHTTPClient(vipr)
			if err != nil {
				logger.WithError(err).Error("failed to initialize Key Manager Client")
				return
			}
			client = qkm.NewHTTPClient(httpClient, &qkm.Config{
				URL: cfg.URL,
			})
			storeNameID = vipr.GetString(StoreNameViperKey)
			logger.WithField("url", cfg.URL).Info("client ready")
		} else {
			client = NewNonClient()
		}
	})
}

// GlobalChainRegistryClient return the chain registry
func GlobalClient() qkm.KeyManagerClient {
	return client
}

func GlobalStoreName() string {
	return storeNameID
}

func SetGlobalStoreName(storeName string) {
	storeNameID = storeName
}
