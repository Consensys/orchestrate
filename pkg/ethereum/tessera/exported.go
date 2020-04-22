package tessera

import (
	"context"
	"sync"

	"github.com/cenkalti/backoff/v4"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

var (
	client   Client
	initOnce = &sync.Once{}
)

func Init(_ context.Context) {
	initOnce.Do(func() {
		if client != nil {
			return
		}

		newBackOff := func() backoff.BackOff { return utils.NewBackOff(utils.NewConfig(viper.GetViper())) }

		client = NewTesseraClient(newBackOff, http.NewClient())
	})
}

// GlobalClient returns global Tessera HttpClient
func GlobalClient() Client {
	return client
}

// SetGlobalClient sets global Tessera HttpClient
func SetGlobalMultiClient(ec Client) {
	client = ec
}
