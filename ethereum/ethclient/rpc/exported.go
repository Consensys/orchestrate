package rpc

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	httpclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/client"
)

const component = "ethclient.rpc"

var (
	client   *Client
	config   *Config
	initOnce = &sync.Once{}
)

func Init(ctx context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		if config == nil {
			config = NewConfig()
		}

		// Set Client
		httpclient.Init(ctx)
		client = NewClient(config, httpclient.GlobalClient())

		log.Infof("%s: ready", component)

	})
}

// GlobalClient returns global Client
func GlobalClient() *Client {
	return client
}

// SetGlobalClient sets global Client
func SetGlobalClient(ec *Client) {
	client = ec
}
