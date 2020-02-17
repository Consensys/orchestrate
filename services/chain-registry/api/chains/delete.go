package chains

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/utils"
	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type deleteResponse struct{}

// @Summary Deletes a chain by ID
// @Produce json
// @Param uuid path string true "ID of the chain"
// @Success 200 {object} deleteResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chains/{uuid} [delete]
func (h Handler) deleteChainByUUID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	err := h.store.DeleteChainByUUID(request.Context(), mux.Vars(request)["uuid"])
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&deleteResponse{})
}

// @Summary Deletes a chain by tenantID and name
// @Produce json
// @Param tenantID path string true "ID of the tenant"
// @Param name path string true "Name of the chain"
// @Success 200 {object} deleteResponse
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /{tenantID}/chains/{name} [delete]
func (h Handler) deleteChainByName(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chain := &models.Chain{
		Name:     mux.Vars(request)["name"],
		TenantID: mux.Vars(request)["tenantID"],
	}

	err := h.store.DeleteChainByName(request.Context(), chain)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&deleteResponse{})
}
