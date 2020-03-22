package chains

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type PatchRequest struct {
	Name     string                `json:"name,omitempty"`
	URLs     []string              `json:"urls,omitempty" pg:"urls,array" validate:"unique,dive,url"`
	Listener *ListenerPatchRequest `json:"listener,omitempty"`
}

type ListenerPatchRequest struct {
	Depth             *uint64 `json:"depth,omitempty"`
	CurrentBlock      *uint64 `json:"currentBlock,string,omitempty"`
	BackOffDuration   *string `json:"backOffDuration,omitempty"`
	ExternalTxEnabled *bool   `json:"externalTxEnabled,omitempty"`
}

// @Summary Updates a chain by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the chain"
// @Param request body PatchRequest true "Chain update request"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chains/{uuid} [patch]
func (h *Chains) PatchChain(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chainRequest := &PatchRequest{Listener: &ListenerPatchRequest{}}
	err := utils.UnmarshalBody(request.Body, chainRequest)
	if err != nil {
		utils.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	uuid := mux.Vars(request)["uuid"]

	chain := &types.Chain{
		UUID:                      uuid,
		Name:                      chainRequest.Name,
		URLs:                      chainRequest.URLs,
		ListenerCurrentBlock:      chainRequest.Listener.CurrentBlock,
		ListenerBackOffDuration:   chainRequest.Listener.BackOffDuration,
		ListenerExternalTxEnabled: chainRequest.Listener.ExternalTxEnabled,
		ListenerDepth:             chainRequest.Listener.Depth,
	}

	err = h.store.UpdateChainByUUID(request.Context(), chain)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	chain, err = h.store.GetChainByUUID(request.Context(), uuid)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chain)
}
