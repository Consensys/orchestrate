package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	healthz "github.com/heptiolabs/healthcheck"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
)

const component = "chain-registry.client"

var (
	client   ChainRegistryClient
	initOnce = &sync.Once{}
	checker  = func() error { return nil }
)

func Init(_ context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		conf := NewConfigFromViper(viper.GetViper())
		client = NewHTTPClient(http.NewClient(http.NewConfig(viper.GetViper())), conf)
		checker = healthz.HTTPGetCheck(fmt.Sprintf("%s/live", viper.GetString(MetricsURLViperKey)), time.Second)
		log.Infof("%s: client ready - url: %s", component, conf.URL)
	})
}

// GlobalChainRegistryClient return the chain registry
func GlobalClient() ChainRegistryClient {
	return client
}

func GlobalChecker() healthz.Check {
	return checker
}

func SetGlobalChecker(c healthz.Check) {
	checker = c
}

// SetGlobalChainRegistryClient set a the chain registry client
func SetGlobalClient(c ChainRegistryClient) {
	client = c
}
