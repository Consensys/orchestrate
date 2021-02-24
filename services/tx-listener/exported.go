package txlistener

import (
	"context"
	"sync"
	"time"

	"github.com/ConsenSys/orchestrate/pkg/app"
	authkey "github.com/ConsenSys/orchestrate/pkg/auth/key"
	authutils "github.com/ConsenSys/orchestrate/pkg/auth/utils"
	"github.com/ConsenSys/orchestrate/pkg/backoff"
	ethclient "github.com/ConsenSys/orchestrate/pkg/ethclient/rpc"
	"github.com/ConsenSys/orchestrate/pkg/http"
	orchestrateclient "github.com/ConsenSys/orchestrate/pkg/sdk/client"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	registryprovider "github.com/ConsenSys/orchestrate/services/tx-listener/providers/chain-registry"
	kafkahook "github.com/ConsenSys/orchestrate/services/tx-listener/session/ethereum/hooks/kafka"
	registryoffset "github.com/ConsenSys/orchestrate/services/tx-listener/session/ethereum/offset/chain-registry"
	"github.com/spf13/viper"
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
