package client

import (
	"context"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

const component = "chain-registry.client"

var (
	client   Client
	initOnce = &sync.Once{}
)

func Init(_ context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		url := viper.GetString(ChainRegistryURLViperKey)
		client = NewHTTPClient(
			http.Client{Timeout: 10 * time.Second},
			Config{URL: url},
		)

		log.Infof("%s: client ready - url: %s", component, url)
	})
}

// GlobalChainRegistryClient return the chain registry
func GlobalClient() Client {
	return client
}

// SetGlobalChainRegistryClient set a the chain registry client
func SetGlobalClient(c Client) {
	client = c
}
