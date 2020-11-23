package utils

import (
	"context"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
)

var chainsCtxKey ctxKey = "chains"

func RegisterNewChain(ctx context.Context, client chainregistry.ChainRegistryClient, chainName string, urls []string) (*models.Chain, error) {
	log.FromContext(ctx).Debugf("Registering new chain '%s' [%q]...", chainName, urls)
	c, err := client.RegisterChain(ctx, &models.Chain{
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

func ContextWithChains(ctx context.Context, chains map[string]*models.Chain) context.Context {
	return context.WithValue(ctx, chainsCtxKey, chains)
}

func ContextChains(ctx context.Context) map[string]*models.Chain {
	v, ok := ctx.Value(chainsCtxKey).(map[string]*models.Chain)
	if !ok {
		return make(map[string]*models.Chain)
	}

	return v
}
