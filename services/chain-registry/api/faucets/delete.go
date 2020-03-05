package faucets

import (
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"

	"github.com/gorilla/mux"
)

// @Summary Deletes a faucet by ID
// @Produce json
// @Param uuid path string true "ID of the faucet"
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /faucets/{uuid} [delete]
func (h Handler) deleteFaucetByUUID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	uuid := mux.Vars(request)["uuid"]

	var err error
	tenantID := multitenancy.TenantIDFromContext(request.Context())
	if tenantID == "" || tenantID == multitenancy.DefaultTenantIDName {
		err = h.store.DeleteFaucetByUUID(request.Context(), uuid)
	} else {
		err = h.store.DeleteFaucetByUUIDAndTenant(request.Context(), uuid, tenantID)
	}

	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
