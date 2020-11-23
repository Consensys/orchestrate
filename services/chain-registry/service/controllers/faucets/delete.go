package faucets

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
)

// @Summary Deletes a faucet by ID
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the faucet"
// @Success 204 "Faucet deleted successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Faucet not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /faucets/{uuid} [delete]
func (h *controller) DeleteFaucet(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	err := h.deleteFaucetUC.Execute(
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
