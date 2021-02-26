package client

import (
	"sync"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/backoff"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/spf13/viper"
)

const component = "api.client"

var (
	client   OrchestrateClient
	initOnce = &sync.Once{}
)

func Init() {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		logger := log.NewLogger().SetComponent(component)

		conf := NewConfigFromViper(viper.GetViper(), backoff.ConstantBackOffWithMaxRetries(time.Second, 5))
		client = NewHTTPClient(http.NewClient(http.NewConfig(viper.GetViper())), conf)
		logger.SetComponent(component).WithField("url", conf.URL).Info("ready")
	})
}

func GlobalClient() OrchestrateClient {
	return client
}
