package faucets

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/utils"
)

// @Summary Deletes a faucet by ID
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the faucet"
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /faucets/{uuid} [delete]
func (h *controller) DeleteFaucet(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	uuid := mux.Vars(request)["uuid"]
	tenantID := multitenancy.TenantIDFromContext(request.Context())

	err := h.deleteFaucet.Execute(request.Context(), uuid, tenantID)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
