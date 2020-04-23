package faucets

import (
	"encoding/json"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/utils"
)

// @Summary Updates a faucet by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the faucet"
// @Param request body PatchRequest true "Faucet update request"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /faucets/{uuid} [patch]
func (h *controller) PatchFaucet(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	faucetRequest := &PatchRequest{}
	err := jsonutils.UnmarshalBody(request.Body, faucetRequest)
	if err != nil {
		utils.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	faucet := parsePatchRequestToFaucet(faucetRequest)
	tenantID := multitenancy.TenantIDFromContext(request.Context())
	if tenantID != "" {
		faucet.TenantID = tenantID
	}

	uuid := mux.Vars(request)["uuid"]
	err = h.updateFaucet.Execute(request.Context(), uuid, faucet)
	if err != nil {
		utils.HandleStoreError(rw, err)
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
