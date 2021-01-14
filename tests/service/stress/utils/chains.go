package utils

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
)

var chainsCtxKey ctxKey = "chains"

func RegisterNewChain(ctx context.Context, client orchestrateclient.OrchestrateClient, chainName string, urls []string) (*api.ChainResponse, error) {
	log.FromContext(ctx).Debugf("Registering new chain '%s' [%q]...", chainName, urls)
	c, err := client.RegisterChain(ctx, &api.RegisterChainRequest{
		Name: chainName,
		URLs: urls,
	})

	if err != nil {
		return nil, err
	}

	if c.UUID == "" {
		return nil, errors.DataError("cannot register chain '%s'", chainName)
	}

	log.FromContext(ctx).Infof("New chain %s registered: %s", chainName, c.UUID)
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
