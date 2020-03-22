package client

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
)

const component = "chain-registry.client"

var (
	client   ChainRegistryClient
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		conf := NewConfig()
		client = NewHTTPClient(
			http.NewClient(),
			conf,
		)

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
