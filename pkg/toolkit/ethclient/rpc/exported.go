package rpc

import (
	"context"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"github.com/consensys/orchestrate/pkg/toolkit/app/auth/key"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/spf13/viper"

	"github.com/consensys/orchestrate/pkg/toolkit/app/http"
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

		vipr := viper.GetViper()
		newBackOff := func() backoff.BackOff { return utils.NewBackOff(utils.NewConfig(vipr)) }

		httpCfg := http.NewDefaultConfig()
		httpCfg.XAPIKey = vipr.GetString(key.APIKeyViperKey)
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
