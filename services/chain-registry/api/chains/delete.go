package chains

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
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
	uuid := mux.Vars(request)["uuid"]

	var err error
	tenantID := multitenancy.TenantIDFromContext(request.Context())
	if tenantID == "" || tenantID == multitenancy.DefaultTenantIDName {
		err = h.store.DeleteChainByUUID(request.Context(), uuid)
	} else {
		err = h.store.DeleteChainByUUIDAndTenant(request.Context(), uuid, tenantID)
	}

	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&deleteResponse{})
}
