package txlistener

import (
	"context"
	"sync"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/app/worker"
	authkey "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/key"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/auth/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/rpc"
	orchlog "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/logger"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	storeclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/envelope-store/client"
	registryprovider "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/providers/chain-registry"
	kafkahook "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/hooks/kafka"
	registryoffset "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/offset/chain-registry"
)

var (
	listener  *TxListener
	appli     *app.App
	initOnce  = &sync.Once{}
	startOnce = &sync.Once{}
	done      chan struct{}
	cancel    func()
)

func initDependencies(ctx context.Context) {
	utils.InParallel(
		func() { registryprovider.Init(ctx) },
		func() { kafkahook.Init(ctx) },
		func() { registryoffset.Init(ctx) },
		func() { rpc.Init(ctx) },
		func() { storeclient.Init(ctx) },
	)
}

// Init hook
func Init(ctx context.Context) {
	initOnce.Do(func() {
		if listener != nil {
			return
		}

		initDependencies(ctx)

		listener = NewTxListener(
			registryprovider.GlobalProvider(),
			kafkahook.GlobalHook(),
			registryoffset.GlobalManager(),
			rpc.GlobalClient(),
			storeclient.GlobalEnvelopeStoreClient(),
		)
	})
}

// Start starts application
func Start(ctx context.Context) error {
	var err error
	startOnce.Do(func() {
		// Create Configuration
		cfg := app.NewConfig(viper.GetViper())
		orchlog.ConfigureLogger(cfg.HTTP)

		ctx, cancel = context.WithCancel(ctx)

		// Create appli to expose metrics
		appli, err = worker.New(cfg, prom.DefaultRegisterer)
		if err != nil {
			return
		}

		err = appli.Start(ctx)
		if err != nil {
			return
		}

		apiKey := viper.GetString(authkey.APIKeyViperKey)
		if apiKey != "" {
			// Inject authorization header in context for later authentication
			ctx = authutils.WithAPIKey(ctx, apiKey)
		}

		Init(ctx)
		done = make(chan struct{})
		go func() {
			listener.Start(ctx)
			close(done)
		}()

	})

	return err
}

func Stop(ctx context.Context) error {
	cancel()
	err := appli.Stop(ctx)
	<-done
	return err
}
