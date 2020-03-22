package chains

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

// @Summary Retrieves a list of all registered chains
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 404
// @Failure 500
// @Router /chains [get]
func (h *Chains) GetChains(rw http.ResponseWriter, request *http.Request) {
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

	if len(chains) == 0 {
		chains = []*types.Chain{}
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
func (h *Chains) GetChain(rw http.ResponseWriter, request *http.Request) {
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
