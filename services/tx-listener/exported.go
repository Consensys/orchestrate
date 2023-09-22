package txlistener

import (
	"context"
	"sync"
	"time"

	"github.com/consensys/orchestrate/pkg/backoff"
	orchestrateclient "github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app"
	authkey "github.com/consensys/orchestrate/pkg/toolkit/app/auth/key"
	authutils "github.com/consensys/orchestrate/pkg/toolkit/app/auth/utils"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	ethclient "github.com/consensys/orchestrate/pkg/toolkit/ethclient/rpc"
	"github.com/consensys/orchestrate/pkg/utils"
	registryprovider "github.com/consensys/orchestrate/services/tx-listener/providers/chain-registry"
	kafkahook "github.com/consensys/orchestrate/services/tx-listener/session/ethereum/hooks/kafka"
	registryoffset "github.com/consensys/orchestrate/services/tx-listener/session/ethereum/offset/chain-registry"
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
		func() {
			kafkahook.Init(ctx)
		},
		func() {
			viper.Set(utils.RetryMaxIntervalViperKey, 30*time.Second)
			viper.Set(utils.RetryMaxElapsedTimeViperKey, 1*time.Hour)
			ethclient.Init(ctx)
		},
	)

	httpClient := app.NewHTTPClient(viper.GetViper())
	backoffConf := orchestrateclient.NewConfigFromViper(viper.GetViper(),
		backoff.IncrementalBackOffWithMaxRetries(time.Millisecond*500, time.Second, 5))
	client := orchestrateclient.NewHTTPClient(httpClient, backoffConf)

	registryprovider.Init(client)
	registryoffset.Init(client)

	//FIXME CUSTOM HEADER
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
		if viper.GetBool(multitenancy.EnabledViperKey) {
			apiKey := viper.GetString(authkey.APIKeyViperKey)
			ctx = multitenancy.WithUserInfo(
				authutils.WithAPIKey(ctx, apiKey),
				multitenancy.NewInternalAdminUser())
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
