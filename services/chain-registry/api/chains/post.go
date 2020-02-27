package chains

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
)

type PostRequest struct {
	Name     string               `json:"name" validate:"required"`
	URLs     []string             `json:"urls" pg:"urls,array" validate:"min=1,unique,dive,url"`
	Listener *ListenerPostRequest `json:"listener,omitempty"`
}

type ListenerPostRequest struct {
	Depth             *uint64 `json:"depth,omitempty"`
	FromBlock         *string `json:"fromBlock,omitempty"`
	BackOffDuration   *string `json:"backOffDuration,omitempty"`
	ExternalTxEnabled *bool   `json:"externalTxEnabled,omitempty"`
}

const LatestBlock string = "latest"

// @Summary Registers a new chain
// @Accept json
// @Produce json
// @Param request body PostRequest true "Chain registration request"
// @Success 200 {object} Chain
// @Failure 400
// @Failure 500
// @Router /chains [post]
func (h Handler) postChain(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chainRequest := &PostRequest{Listener: &ListenerPostRequest{}}
	err := utils.UnmarshalBody(request.Body, chainRequest)
	if err != nil {
		utils.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	chain, err := h.newChain(request.Context(), chainRequest)
	if err != nil {
		utils.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.store.RegisterChain(request.Context(), chain)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chain)
}

func (h Handler) newChain(ctx context.Context, request *PostRequest) (*Chain, error) {
	startingBlock, err := h.processStartingBlock(ctx, request.Listener.FromBlock, request.URLs)
	if err != nil {
		return nil, err
	}

	chain := &Chain{
		Name:                      request.Name,
		URLs:                      request.URLs,
		TenantID:                  multitenancy.TenantIDFromContext(ctx),
		ListenerStartingBlock:     &startingBlock,
		ListenerBackOffDuration:   request.Listener.BackOffDuration,
		ListenerDepth:             request.Listener.Depth,
		ListenerExternalTxEnabled: request.Listener.ExternalTxEnabled,
	}
	chain.SetDefault()

	return chain, nil
}

func (h Handler) processStartingBlock(ctx context.Context, fromBlock *string, urls []string) (uint64, error) {
	if fromBlock == nil || *fromBlock == "" || *fromBlock == LatestBlock {
		return h.getChainTip(ctx, urls)
	}

	return strconv.ParseUint(*fromBlock, 10, 64)
}
