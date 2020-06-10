package chains

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/utils"
)

// @Summary Deletes a chain by ID
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the chain"
// @Success 204
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chains/{uuid} [delete]
func (h *controller) DeleteChain(rw http.ResponseWriter, request *http.Request) {
	err := h.deleteChainUC.Execute(
		request.Context(),
		mux.Vars(request)["uuid"],
		multitenancy.AllowedTenantsFromContext(request.Context()),
	)

	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
