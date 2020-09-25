package txlistener

import (
	"context"
	"sync"
	"time"

	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/backoff"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"
	registryprovider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/providers/chain-registry"
	kafkahook "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/hooks/kafka"
	registryoffset "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/offset/chain-registry"
	txsentry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-sentry"
)

var (
	listener  *TxListener
	appli     *app.App
	sentry    app.Daemon
	initOnce  = &sync.Once{}
	startOnce = &sync.Once{}
)

func initDependencies(ctx context.Context) {
	utils.InParallel(
		func() { registryprovider.Init(ctx) },
		func() { kafkahook.Init(ctx) },
		func() { registryoffset.Init(ctx) },
		func() { ethclient.Init(ctx) },
	)
}

// Init hook
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if listener != nil {
			return
		}

		initDependencies(ctx)

		httpClient := http.NewClient(http.NewConfig(viper.GetViper()))
		backoffConf := txscheduler.NewConfigFromViper(viper.GetViper(), backoff.ConstantBackOffWithMaxRetries(time.Second, 5))
		txSchedulerClientListener := txscheduler.NewHTTPClient(httpClient, backoffConf)
		listener = NewTxListener(
			registryprovider.GlobalProvider(),
			kafkahook.GlobalHook(),
			registryoffset.GlobalManager(),
			ethclient.GlobalClient(),
			txSchedulerClientListener,
		)

		conf := txscheduler.NewConfigFromViper(viper.GetViper(), nil)
		txSchedulerClientSentry := txscheduler.NewHTTPClient(httpClient, conf)
		sentry = txsentry.NewTxSentry(
			txSchedulerClientSentry,
			txsentry.NewConfig(viper.GetViper()),
		)
	})
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

		Init(ctx)

		// Create appli
		appli, err = New(
			app.NewConfig(viper.GetViper()),
			listener,
			sentry,
		)
		if err != nil {

			return
		}

		err = appli.Run(ctx)
	})

	return err
}
