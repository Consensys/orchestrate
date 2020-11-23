package chains

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
)

// @Summary Retrieves a list of all registered chains
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200 {array} models.Chain
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /chains [get]
func (h *controller) GetChains(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	filters := httputil.ToFilters(request.URL.Query())
	chains, err := h.getChainsUC.Execute(
		request.Context(),
		multitenancy.AllowedTenantsFromContext(request.Context()),
		filters,
	)

	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
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
// @Success 200 {object} models.Chain
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Chain not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /chains/{uuid} [get]
func (h *controller) GetChain(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chain, err := h.getChainUC.Execute(
		request.Context(),
		mux.Vars(request)["uuid"],
		multitenancy.AllowedTenantsFromContext(request.Context()),
	)

	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chain)
}
