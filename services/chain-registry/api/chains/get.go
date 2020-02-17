package chains

import (
	"encoding/json"
	"html"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/utils"
)

// @Summary Retrieves a list of all registered chains
// @Produce json
// @Success 200
// @Failure 404
// @Failure 500
// @Router /chains [get]
func (h Handler) getChains(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	filters := make(map[string]string)
	for k := range request.URL.Query() {
		key := html.EscapeString(k)
		filters[key] = html.EscapeString(request.URL.Query().Get(k))
	}

	chains, err := h.store.GetChains(request.Context(), filters)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chains)
}

// @Summary Retrieves a chain by ID
// @Produce json
// @Param uuid path string true "ID of the chain"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chains/{uuid} [get]
func (h Handler) getChainByUUID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chain, err := h.store.GetChainByUUID(request.Context(), mux.Vars(request)["uuid"])
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chain)
}

// @Summary Retrieves a list of all registered chains by tenantID
// @Produce json
// @Param tenantID path string true "ID of the tenant"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /{tenantID}/chains [get]
func (h Handler) getChainsByTenantID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	filters := make(map[string]string)
	for k := range request.URL.Query() {
		key := html.EscapeString(k)
		filters[key] = html.EscapeString(request.URL.Query().Get(k))
	}

	chains, err := h.store.GetChainsByTenantID(request.Context(), mux.Vars(request)["tenantID"], filters)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chains)
}

// @Summary Retrieves a chain by tenantID and name
// @Produce json
// @Param tenantID path string true "ID of the tenant"
// @Param name path string true "Name of the chain"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /{tenantID}/chains/{name} [get]
func (h Handler) getChainByTenantIDAndName(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chain, err := h.store.GetChainByTenantIDAndName(request.Context(), mux.Vars(request)["tenantID"], mux.Vars(request)["name"])
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chain)
}
