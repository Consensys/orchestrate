package chains

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
)

// @Summary Retrieves a list of all registered chains
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200 {array} Chain
// @Failure 404
// @Failure 500
// @Router /chains [get]
func (h Handler) getChains(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	filters := utils.ToFilters(request.URL.Query())
	var chains []*types.Chain
	var err error
	tenantID := multitenancy.TenantIDFromContext(request.Context())
	if tenantID == "" || tenantID == multitenancy.DefaultTenantIDName {
		chains, err = h.store.GetChains(request.Context(), filters)
	} else {
		chains, err = h.store.GetChainsByTenant(request.Context(), filters, tenantID)
	}

	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chains)
}

// @Summary Retrieves a chain by ID
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the chain"
// @Success 200 {object} Chain
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chains/{uuid} [get]
func (h Handler) getChainByUUID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	uuid := mux.Vars(request)["uuid"]

	var chain *types.Chain
	var err error
	tenantID := multitenancy.TenantIDFromContext(request.Context())
	if tenantID == "" || tenantID == multitenancy.DefaultTenantIDName {
		chain, err = h.store.GetChainByUUID(request.Context(), uuid)
	} else {
		chain, err = h.store.GetChainByUUIDAndTenant(request.Context(), uuid, tenantID)
	}

	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chain)
}
