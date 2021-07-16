package quorumkeymanager

import (
	"sync"

	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
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
	cfg := NewConfigFromViper(viper.GetViper())
	initOnce.Do(func() {
		vipr := viper.GetViper()
		if client != nil {
			return
		}
		logger := log.NewLogger().SetComponent(component)
		client = qkm.NewHTTPClient(http.NewClient(http.NewConfig(vipr)), &qkm.Config{
			URL: cfg.URL,
		})
		storeNameID = vipr.GetString(StoreNameViperKey)
		logger.WithField("url", cfg.URL).Info("client ready")
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
