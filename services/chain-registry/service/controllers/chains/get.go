package chains

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/utils"
)

// @Summary Retrieves a list of all registered chains
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 404
// @Failure 500
// @Router /chains [get]
func (h *controller) GetChains(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	filters := utils.ToFilters(request.URL.Query())
	tenantID := multitenancy.TenantIDFromContext(request.Context())
	chains, err := h.getChainsUC.Execute(request.Context(), tenantID, filters)

	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	if len(chains) == 0 {
		chains = []*models.Chain{}
	}

	_ = json.NewEncoder(rw).Encode(chains)
}

// @Summary Retrieves a chain by ID
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the chain"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chains/{uuid} [get]
func (h *controller) GetChain(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	uuid := mux.Vars(request)["uuid"]

	tenantID := multitenancy.TenantIDFromContext(request.Context())
	chain, err := h.getChainUC.Execute(request.Context(), uuid, tenantID)

	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chain)
}
