package utils

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	utils4 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"

	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
)

var chainsCtxKey ctxKey = "chains"

func RegisterNewChain(
	ctx context.Context,
	client orchestrateclient.OrchestrateClient,
	ec ethclient.Client,
	chainName string,
	urls []string,
) (*api.ChainResponse, error) {
	log.WithContext(ctx).Debugf("Registering new chain '%s' [%q]...", chainName, urls)
	c, err := client.RegisterChain(ctx, &api.RegisterChainRequest{
		Name: chainName,
		URLs: urls,
	})
	if err != nil {
		return nil, err
	}

	// Give time to the proxy to be set
	apiURL := viper.GetString(orchestrateclient.URLViperKey)
	proxyURL := utils4.GetProxyURL(apiURL, c.UUID)
	err = backoff.RetryNotify(
		func() error {
			_, err2 := ec.Network(ctx, proxyURL)
			return err2
		},
		backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Second), 5),
		func(err error, duration time.Duration) {
			log.WithContext(ctx).WithField("chain_uuid", c.UUID).WithError(err).Debug("scenario: chain proxy is still not ready")
		},
	)
	if err != nil {
		return nil, err
	}

	log.WithContext(ctx).Infof("New chain %s registered: %s", chainName, c.UUID)
	return c, nil
}

func ContextWithChains(ctx context.Context, chains map[string]*api.ChainResponse) context.Context {
	return context.WithValue(ctx, chainsCtxKey, chains)
}

func ContextChains(ctx context.Context) map[string]*api.ChainResponse {
	v, ok := ctx.Value(chainsCtxKey).(map[string]*api.ChainResponse)
	if !ok {
		return make(map[string]*api.ChainResponse)
	}

	return v
}
