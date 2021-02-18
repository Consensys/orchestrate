package txlistener

import (
	"context"
	"sync"
	"time"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/app"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/key"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/backoff"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	registryprovider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/providers/chain-registry"
	kafkahook "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/session/ethereum/hooks/kafka"
	registryoffset "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/tx-listener/session/ethereum/offset/chain-registry"
)

var (
	listener  *TxListener
	appli     *app.App
	sentry    app.Daemon
	startOnce = &sync.Once{}
)

// New Utility function used to initialize a new service
func NewApp(ctx context.Context) (*app.App, error) {
	// Initialize dependencies
	config := app.NewConfig(viper.GetViper())

	utils.InParallel(
		func() { kafkahook.Init(ctx) },
		func() {
			viper.Set(utils.RetryMaxIntervalViperKey, 30*time.Second)
			viper.Set(utils.RetryMaxElapsedTimeViperKey, 1*time.Hour)
			ethclient.Init(ctx)
		},
	)

	httpClient := http.NewClient(http.NewConfig(viper.GetViper()))
	backoffConf := orchestrateclient.NewConfigFromViper(viper.GetViper(),
		backoff.IncrementalBackOffWithMaxRetries(time.Millisecond*500, time.Second, 5))
	client := orchestrateclient.NewHTTPClient(httpClient, backoffConf)

	registryprovider.Init(client)
	registryoffset.Init(client)

	return New(
		config,
		registryprovider.GlobalProvider(),
		kafkahook.GlobalHook(),
		registryoffset.GlobalManager(),
		ethclient.GlobalClient(),
		client,
	)
}

// Start starts application
func Run(ctx context.Context) error {
	var err error
	startOnce.Do(func() {
		apiKey := viper.GetString(authkey.APIKeyViperKey)
		if apiKey != "" {
			// Inject authorization header in context for later authentication
			ctx = authutils.WithAPIKey(ctx, apiKey)
		}

		// Create appli
		appli, err = NewApp(ctx)
		if err != nil {

			return
		}

		err = appli.Run(ctx)
	})

	return err
}
