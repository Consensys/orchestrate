package formatters

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/consensys/orchestrate/pkg/errors"
	log "github.com/sirupsen/logrus"

	types "github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
)

func FormatChainResponse(chain *entities.Chain) *types.ChainResponse {
	return &types.ChainResponse{
		UUID:                      chain.UUID,
		Name:                      chain.Name,
		TenantID:                  chain.TenantID,
		OwnerID:                   chain.OwnerID,
		URLs:                      chain.URLs,
		ChainID:                   chain.ChainID,
		ListenerDepth:             chain.ListenerDepth,
		ListenerCurrentBlock:      chain.ListenerCurrentBlock,
		ListenerStartingBlock:     chain.ListenerStartingBlock,
		ListenerBackOffDuration:   chain.ListenerBackOffDuration,
		ListenerExternalTxEnabled: chain.ListenerExternalTxEnabled,
		PrivateTxManager:          chain.PrivateTxManager,
		Labels:                    chain.Labels,
		Headers:                   chain.Headers,
		CreatedAt:                 chain.CreatedAt,
		UpdatedAt:                 chain.UpdatedAt,
	}
}

func FormatRegisterChainRequest(request *types.RegisterChainRequest, fromLatest bool) (*entities.Chain, error) {
	chain := &entities.Chain{
		Name:                      request.Name,
		URLs:                      request.URLs,
		ListenerDepth:             request.Listener.Depth,
		ListenerBackOffDuration:   request.Listener.BackOffDuration,
		ListenerExternalTxEnabled: request.Listener.ExternalTxEnabled,
		Labels:                    request.Labels,
		Headers:                   request.Headers,
	}

	if request.Listener.BackOffDuration == "" {
		chain.ListenerBackOffDuration = "5s"
	}

	if !fromLatest {
		startingBlock, err := strconv.ParseUint(request.Listener.FromBlock, 10, 64)
		if err != nil {
			errMessage := "fromBlock must be an integer value"
			log.WithField("from_block", request.Listener.FromBlock).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}
		chain.ListenerStartingBlock = startingBlock
		chain.ListenerCurrentBlock = startingBlock
	}

	if request.PrivateTxManager != nil {
		chain.PrivateTxManager = &entities.PrivateTxManager{
			URL:  request.PrivateTxManager.URL,
			Type: request.PrivateTxManager.Type,
		}
	}

	return chain, nil
}

func FormatUpdateChainRequest(request *types.UpdateChainRequest, uuid string) *entities.Chain {
	chain := &entities.Chain{
		UUID:    uuid,
		Name:    request.Name,
		Labels:  request.Labels,
		Headers: request.Headers,
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
