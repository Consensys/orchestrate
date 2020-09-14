package chains

import (
	"encoding/json"
	"net/http"
	"strconv"

	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

const latestBlockStr string = "latest"

// @Summary Registers a new chain
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body PostRequest true "Chain registration request"
// @Success 200 {object} models.Chain
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /chains [post]
func (h *controller) PostChain(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chain, err := parsePostReqToChain(request)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.registerChainUC.Execute(request.Context(), chain)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chain)
}

func parsePostReqToChain(request *http.Request) (*models.Chain, error) {
	chainRequest := &PostRequest{
		Listener: &ListenerPostRequest{},
	}

	err := jsonutils.UnmarshalBody(request.Body, chainRequest)
	if err != nil {
		return nil, err
	}

	var listenerStartingBlock *uint64
	fromBlock := chainRequest.Listener.FromBlock
	if fromBlock != nil && *fromBlock != "" && *fromBlock != latestBlockStr {
		head, err := strconv.ParseUint(*fromBlock, 10, 64)
		if err != nil {
			return nil, err
		}

		listenerStartingBlock = &head
	}

	chain := models.Chain{
		Name:                      chainRequest.Name,
		URLs:                      chainRequest.URLs,
		TenantID:                  multitenancy.TenantIDFromContext(request.Context()),
		ListenerStartingBlock:     listenerStartingBlock,
		ListenerBackOffDuration:   chainRequest.Listener.BackOffDuration,
		ListenerDepth:             chainRequest.Listener.Depth,
		ListenerExternalTxEnabled: chainRequest.Listener.ExternalTxEnabled,
	}

	if chainRequest.PrivateTxManager != nil {
		chain.PrivateTxManagers = []*models.PrivateTxManagerModel{
			{
				URL:  chainRequest.PrivateTxManager.URL,
				Type: chainRequest.PrivateTxManager.Type,
			},
		}
	}

	return &chain, nil
}
