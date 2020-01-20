package client

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	httpclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client"
)

const component = "chain-registry.client"

var (
	client   Client
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		httpclient.Init(ctx)

		url := viper.GetString(ChainRegistryURLViperKey)
		client = NewHTTPClient(
			httpclient.GlobalClient(),
			&Config{URL: url},
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
