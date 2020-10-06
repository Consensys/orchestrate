package client

import (
	"fmt"
	"sync"
	"time"

	healthz "github.com/heptiolabs/healthcheck"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/backoff"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
)

const component = "transaction-scheduler.client"

var (
	client   TransactionSchedulerClient
	initOnce = &sync.Once{}
	checker  healthz.Check
)

func Init() {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		conf := NewConfigFromViper(viper.GetViper(), backoff.ConstantBackOffWithMaxRetries(time.Second, 5))
		client = NewHTTPClient(http.NewClient(http.NewConfig(viper.GetViper())), conf)
		checker = healthz.HTTPGetCheck(fmt.Sprintf("%s/live", viper.GetString(MetricsURLViperKey)), time.Second)
		log.Infof("%s: client ready - url: %s", component, conf.URL)
	})
}

// GlobalChainRegistryClient return the chain registry
func GlobalClient() TransactionSchedulerClient {
	return client
}

func GlobalChecker() healthz.Check {
	return checker
}

func SetGlobalChecker(c healthz.Check) {
	checker = c
}
