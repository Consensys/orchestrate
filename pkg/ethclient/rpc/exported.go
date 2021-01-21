package rpc

import (
	"context"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
)

const component = "ethclient.rpc"

var (
	client   *Client
	initOnce = &sync.Once{}
)

func Init(_ context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}
		logger := log.NewLogger().SetComponent(component)

		newBackOff := func() backoff.BackOff { return utils.NewBackOff(utils.NewConfig(viper.GetViper())) }

		httpCfg := http.NewConfig(viper.GetViper())

		// Deactivate context authToken forwarding for RPC client requests
		httpCfg.AuthHeaderForward = false

		// Set Client
		client = NewClient(newBackOff, http.NewClient(httpCfg))

		logger.Info("ready")
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
