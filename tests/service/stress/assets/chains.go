package assets

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	utils2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/service/stress/utils"
)

var chainsCtxKey ctxKey = "chains"

type Chain struct {
	Name            string
	ProxyURL        string
	PrivNodeAddress []string
}

func RegisterNewChain(ctx context.Context, client orchestrateclient.OrchestrateClient, ec ethclient.Client,
	proxyHost, chainName string, chainData *utils2.TestDataChain,
) (context.Context, error) {
	logger := log.FromContext(ctx).WithField("name", chainName).WithField("urls", chainData.URLs)
	logger.WithContext(ctx).Debug("registering new chain")

	c, err := client.RegisterChain(ctx, &api.RegisterChainRequest{
		Name: chainName,
		URLs: chainData.URLs,
	})
	if err != nil {
		errMsg := "failed to register chain"
		logger.WithError(err).Error("failed to register chain")
		return nil, fmt.Errorf(errMsg)
	}

	chainProxyURL := utils.GetProxyURL(proxyHost, c.UUID)
	err = backoff.RetryNotify(
		func() error {
			_, err2 := ec.Network(ctx, chainProxyURL)
			return err2
		},
		backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Second), 5),
		func(err error, duration time.Duration) {
			logger.WithField("chain", c.UUID).WithError(err).Debug("chain proxy is still not ready")
		},
	)
	if err != nil {
		return nil, err
	}

	logger.WithField("chain", c.UUID).Info("new chain has been registered")
	return contextWithChains(ctx, append(ContextChains(ctx),
		Chain{ProxyURL: chainProxyURL, Name: chainName, PrivNodeAddress: chainData.PrivateAddress}),
	), nil
}

func contextWithChains(ctx context.Context, chains []Chain) context.Context {
	return context.WithValue(ctx, chainsCtxKey, chains)
}

func ContextChains(ctx context.Context) []Chain {
	if v, ok := ctx.Value(chainsCtxKey).([]Chain); ok {
		return v
	}

	return make([]Chain, 0)
}
