package assets

import (
	"context"
	"fmt"

	orchestrateclient "github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/ethclient"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/utils"
	utils2 "github.com/consensys/orchestrate/tests/utils"
)

var chainsCtxKey ctxKey = "chains"

type Chain struct {
	UUID            string
	Name            string
	ProxyURL        string
	PrivNodeAddress []string
}

func RegisterNewChain(ctx context.Context, client orchestrateclient.OrchestrateClient, ec ethclient.ChainSyncReader,
	proxyHost, chainName string, chainData *utils2.TestDataChain,
) (context.Context, string, error) {
	logger := log.FromContext(ctx).WithField("name", chainName).WithField("urls", chainData.URLs)
	logger.WithContext(ctx).Debug("registering new chain")

	c, err := client.RegisterChain(ctx, &api.RegisterChainRequest{
		Name: chainName,
		URLs: chainData.URLs,
	})
	if err != nil {
		logger.WithError(err).Error("failed to register chain")
		return nil, "", err
	}

	err = utils2.WaitForProxy(ctx, proxyHost, c.UUID, ec)
	if err != nil {
		logger.WithError(err).Error("failed to wait for chain proxy")
		return nil, "", err
	}

	logger.WithField("chain", c.UUID).Info("new chain has been registered")
	return contextWithChains(ctx, append(ContextChains(ctx),
		Chain{
			UUID:            c.UUID,
			ProxyURL:        utils.GetProxyURL(proxyHost, c.UUID),
			Name:            chainName,
			PrivNodeAddress: chainData.PrivateAddress,
		}),
	), c.UUID, nil
}

func DeregisterChain(ctx context.Context, client orchestrateclient.OrchestrateClient, chain *Chain) error {
	logger := log.FromContext(ctx).WithField("uuid", chain.UUID).WithField("name", chain.Name)
	logger.WithContext(ctx).Debug("deleting chain")

	err := client.DeleteChain(ctx, chain.UUID)
	if err != nil {
		errMsg := "failed to delete chain"
		logger.WithError(err).Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	logger.Info("chain has been deleted successfully")
	return nil
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
