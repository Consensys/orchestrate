package client

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/backoff"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
)

const component = "transaction-scheduler.client"

var (
	client   TransactionSchedulerClient
	initOnce = &sync.Once{}
)

func Init() {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		conf := NewConfigFromViper(viper.GetViper(), backoff.ConstantBackOffWithMaxRetries(time.Second, 5))
		client = NewHTTPClient(http.NewClient(), conf)

		log.Infof("%s: client ready - url: %s", component, conf.URL)
	})
}

// GlobalChainRegistryClient return the chain registry
func GlobalClient() TransactionSchedulerClient {
	return client
}
