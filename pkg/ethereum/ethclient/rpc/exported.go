package rpc

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
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
		client = NewClient(config, http.NewClient())

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
