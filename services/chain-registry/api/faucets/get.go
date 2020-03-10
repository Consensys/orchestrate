package faucets

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/utils"
	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
)

// @Summary Retrieves a list of all registered faucet
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200 {array} Faucet
// @Failure 404
// @Failure 500
// @Router /faucets [get]
func (h Handler) getFaucets(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	filters := utils.ToFilters(request.URL.Query())
	var faucets []*models.Faucet
	var err error
	tenantID := multitenancy.TenantIDFromContext(request.Context())
	if tenantID == "" || tenantID == multitenancy.DefaultTenantIDName {
		faucets, err = h.store.GetFaucets(request.Context(), filters)
	} else {
		faucets, err = h.store.GetFaucetsByTenant(request.Context(), filters, tenantID)
	}

	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(faucets)
}

// @Summary Retrieves a faucet by ID
// @Produce json
// @Param uuid path string true "ID of the faucet"
// @Success 200 {object} Faucet
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /faucets/{uuid} [get]
func (h Handler) getFaucetByUUID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	uuid := mux.Vars(request)["uuid"]

	var faucet *models.Faucet
	var err error
	tenantID := multitenancy.TenantIDFromContext(request.Context())
	if tenantID == "" || tenantID == multitenancy.DefaultTenantIDName {
		faucet, err = h.store.GetFaucetByUUID(request.Context(), uuid)
	} else {
		faucet, err = h.store.GetFaucetByUUIDAndTenant(request.Context(), uuid, tenantID)
	}

	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(faucet)
}
