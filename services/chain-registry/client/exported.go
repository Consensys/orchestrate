package client

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
)

const component = "chain-registry.client"

var (
	client   ChainRegistryClient
	initOnce = &sync.Once{}
)

func Init(_ context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		conf := NewConfigFromViper(viper.GetViper())
		client = NewHTTPClient(http.NewClient(http.NewConfig(viper.GetViper())), conf)
		log.Infof("%s: client ready - url: %s", component, conf.URL)
	})
}

// GlobalChainRegistryClient return the chain registry
func GlobalClient() ChainRegistryClient {
	return client
}

// SetGlobalChainRegistryClient set a the chain registry client
func SetGlobalClient(c ChainRegistryClient) {
	client = c
}
