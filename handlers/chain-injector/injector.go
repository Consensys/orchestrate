package chaininjector

import (
	"fmt"

	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/proxy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/engine"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
)

// NodeInjector enrich the envelope with the nodeID, nodeName and inject in the input.Context the proxy URL
func NodeInjector(multitenancyEnabled bool, r registry.Client, chainRegistryURL string) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		// check nodeID and inject if not present
		tenantID := multitenancy.DefaultTenantIDName
		if multitenancyEnabled {
			if tenantID = multitenancy.TenantIDFromContext(txctx.Context()); tenantID == "" {
				e := txctx.AbortWithError(errors.InternalError("invalid tenantID not found")).ExtendComponent(component)
				txctx.Logger.Error(e)
				return
			}
		}

		// Check if node exist
		node, err := getNode(txctx, r, tenantID)
		if err != nil {
			e := txctx.AbortWithError(errors.FromError(err)).ExtendComponent(component)
			txctx.Logger.Error(e)
			return
		}

		// Re-inject node Name and node ID from chain registry
		txctx.Envelope.GetChain().SetNodeName(node.Name)
		txctx.Envelope.GetChain().SetNodeID(node.ID)

		// Inject chain proxy path as /tenantID/nodeName
		proxyURL := fmt.Sprintf("%s/%s", chainRegistryURL, proxy.PathByNodeName(tenantID, node.Name))
		txctx.WithContext(proxy.With(txctx.Context(), proxyURL))
	}
}

// getNode retrieves node in the chain registry
func getNode(txctx *engine.TxContext, r registry.Client, tenantID string) (*types.Node, error) {
	var n *types.Node
	var err error
	nodeID := txctx.Envelope.GetChain().GetNodeId()
	nodeName := txctx.Envelope.GetChain().GetNodeName()
	switch {
	case nodeID != "":
		n, err = r.GetNodeByTenantAndNodeID(txctx.Context(), tenantID, nodeID)
	case nodeName != "":
		n, err = r.GetNodeByTenantAndNodeName(txctx.Context(), tenantID, nodeName)
	default:
		return nil, errors.InternalError("invalid envelope - no node id or node name are filled - cannot retrieve chain id")
	}
	if err != nil {
		return nil, err
	}
	return n, nil
}

// ChainIDInjector enrich the envelope with the chain ID retrieved from the chain proxy
func ChainIDInjector(ec ethclient.ChainSyncReader) engine.HandlerFunc {
	return func(txctx *engine.TxContext) {
		chainProxyURL := proxy.FromContext(txctx.Context())
		chainID, err := ec.Network(txctx.Context(), chainProxyURL)
		if err != nil {
			e := txctx.AbortWithError(err).ExtendComponent(component)
			txctx.Logger.WithError(e).Errorf("injector: could not retrieve chain id from %s", chainProxyURL)
			return
		}
		txctx.Envelope.GetChain().SetID(chainID)
		txctx.Logger.Debugf("injector: chain id %s injected", chainID.String())
	}
}
