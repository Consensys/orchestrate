package chains

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/utils"
	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
)

type PatchRequest struct {
	Name     string    `json:"name,omitempty"`
	URLs     []string  `json:"urls,omitempty" pg:"urls,array" validate:"unique,dive,url"`
	Listener *Listener `json:"listener,omitempty"`
}

type PatchResponse struct{}

// @Summary Updates a chain by ID
// @Accept json
// @Produce json
// @Param uuid path string true "ID of the chain"
// @Param request body PatchRequest true "Chain update request"
// @Success 200 {object} PatchResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chains/{uuid} [patch]
func (h Handler) patchChainByUUID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chainRequest := &PatchRequest{Listener: &Listener{}}
	err := utils.UnmarshalBody(request.Body, chainRequest)
	if err != nil {
		utils.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	chain := &models.Chain{
		UUID:                      mux.Vars(request)["uuid"],
		Name:                      chainRequest.Name,
		URLs:                      chainRequest.URLs,
		ListenerDepth:             chainRequest.Listener.Depth,
		ListenerBlockPosition:     chainRequest.Listener.BlockPosition,
		ListenerBackOffDuration:   chainRequest.Listener.BackOffDuration,
		ListenerExternalTxEnabled: chainRequest.Listener.ExternalTxEnabled,
	}

	tenantID := multitenancy.TenantIDFromContext(request.Context())
	if tenantID != "" {
		chain.TenantID = tenantID
	}

	err = h.store.UpdateChainByUUID(request.Context(), chain)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&PatchResponse{})
}
