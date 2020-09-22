package rpc

import (
	"context"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
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

		newBackOff := func() backoff.BackOff { return utils.NewBackOff(utils.NewConfig(viper.GetViper())) }
		// Set Client
		client = NewClient(newBackOff, http.NewClient(http.NewConfig(viper.GetViper())))

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
