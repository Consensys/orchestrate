package client

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
)

const component = "transaction-scheduler.client"

var (
	client   TransactionSchedulerClient
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		conf := NewConfigFromViper(viper.GetViper())
		client = NewHTTPClient(http.NewClient(), conf)

		log.Infof("%s: client ready - url: %s", component, conf.URL)
	})
}

// GlobalChainRegistryClient return the chain registry
func GlobalClient() TransactionSchedulerClient {
	return client
}

// SetGlobalChainRegistryClient set a the chain registry client
func SetGlobalClient(c TransactionSchedulerClient) {
	client = c
}
