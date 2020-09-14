package chains

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
)

// @Summary Deletes a chain by ID
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the chain"
// @Success 204 "Chain deleted successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Chain not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /chains/{uuid} [delete]
func (h *controller) DeleteChain(rw http.ResponseWriter, request *http.Request) {
	err := h.deleteChainUC.Execute(
		request.Context(),
		mux.Vars(request)["uuid"],
		multitenancy.AllowedTenantsFromContext(request.Context()),
	)

	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
