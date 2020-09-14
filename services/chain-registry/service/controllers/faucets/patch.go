package faucets

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

// @Summary Updates a faucet by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the faucet"
// @Param request body PatchRequest true "Faucet update request"
// @Success 200 {object} models.Faucet
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Chain not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /faucets/{uuid} [patch]
func (h *controller) PatchFaucet(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	faucetRequest := &PatchRequest{}
	err := jsonutils.UnmarshalBody(request.Body, faucetRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	faucet := parsePatchRequestToFaucet(faucetRequest)

	err = h.updateFaucetUC.Execute(
		request.Context(),
		mux.Vars(request)["uuid"],
		multitenancy.AllowedTenantsFromContext(request.Context()),
		faucet)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(rw).Encode(faucet)
}

func parsePatchRequestToFaucet(faucetRequest *PatchRequest) *models.Faucet {
	return &models.Faucet{
		Name:            faucetRequest.Name,
		ChainRule:       faucetRequest.ChainRule,
		CreditorAccount: faucetRequest.CreditorAccount,
		MaxBalance:      faucetRequest.MaxBalance,
		Amount:          faucetRequest.Amount,
		Cooldown:        faucetRequest.Cooldown,
	}
}
