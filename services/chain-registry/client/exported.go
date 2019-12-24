package client

import (
	"context"
	"net/http"
	"sync"
	"time"

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

		client = &HTTPClient{
			client: http.Client{
				Timeout: 10 * time.Second,
			},
			config: Config{
				url: viper.GetString(ChainRegistryURLViperKey)},
		}
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
