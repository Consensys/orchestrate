package formatters

import (
	"net/http"
	"strings"

	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

func FormatChainResponse(chain *entities.Chain) *types.ChainResponse {
	return &types.ChainResponse{
		UUID:                      chain.UUID,
		Name:                      chain.Name,
		TenantID:                  chain.TenantID,
		URLs:                      chain.URLs,
		ChainID:                   chain.ChainID,
		ListenerDepth:             chain.ListenerDepth,
		ListenerCurrentBlock:      chain.ListenerCurrentBlock,
		ListenerStartingBlock:     chain.ListenerStartingBlock,
		ListenerBackOffDuration:   chain.ListenerBackOffDuration,
		ListenerExternalTxEnabled: chain.ListenerExternalTxEnabled,
		PrivateTxManager:          chain.PrivateTxManager,
		CreatedAt:                 chain.CreatedAt,
		UpdatedAt:                 chain.UpdatedAt,
	}
}

func FormatRegisterChainRequest(request *types.RegisterChainRequest, tenantID string) *entities.Chain {
	chain := &entities.Chain{
		Name:                      request.Name,
		TenantID:                  tenantID,
		URLs:                      request.URLs,
		ListenerDepth:             request.Listener.Depth,
		ListenerBackOffDuration:   request.Listener.BackOffDuration,
		ListenerExternalTxEnabled: request.Listener.ExternalTxEnabled,
	}

	if request.Listener.BackOffDuration == "" {
		chain.ListenerBackOffDuration = "5s"
	}

	if request.PrivateTxManager != nil {
		chain.PrivateTxManager = &entities.PrivateTxManager{
			URL:  request.PrivateTxManager.URL,
			Type: request.PrivateTxManager.Type,
		}
	}

	return chain
}

func FormatUpdateChainRequest(request *types.UpdateChainRequest, uuid string) *entities.Chain {
	chain := &entities.Chain{
		UUID: uuid,
		Name: request.Name,
	}

	if request.Listener != nil {
		chain.ListenerBackOffDuration = request.Listener.BackOffDuration
		chain.ListenerDepth = request.Listener.Depth
		chain.ListenerExternalTxEnabled = request.Listener.ExternalTxEnabled
		chain.ListenerCurrentBlock = request.Listener.CurrentBlock
	}

	if request.PrivateTxManager != nil {
		chain.PrivateTxManager = &entities.PrivateTxManager{
			URL:  request.PrivateTxManager.URL,
			Type: request.PrivateTxManager.Type,
		}
	}

	return chain
}

func FormatChainFiltersRequest(req *http.Request) (*entities.ChainFilters, error) {
	filters := &entities.ChainFilters{}

	qNames := req.URL.Query().Get("names")
	if qNames != "" {
		filters.Names = strings.Split(qNames, ",")
	}

	if err := utils.GetValidator().Struct(filters); err != nil {
		return nil, err
	}

	return filters, nil
}
