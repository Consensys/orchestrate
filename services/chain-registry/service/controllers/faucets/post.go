package faucets

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/utils"
)

// @Summary Registers a new faucet
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body PostRequest true "Faucet registration request"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /faucets [post]
func (h *controller) PostFaucet(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	faucetRequest := &PostRequest{}

	err := utils.UnmarshalBody(request.Body, faucetRequest)
	if err != nil {
		utils.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	tenantID := multitenancy.TenantIDFromContext(request.Context())
	if tenantID == "" {
		tenantID = multitenancy.DefaultTenantIDName
	}

	faucet := parsePostRequestToFaucet(faucetRequest)
	faucet.TenantID = tenantID

	err = h.registerFaucet.Execute(request.Context(), faucet)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(faucet)
}

func parsePostRequestToFaucet(request *PostRequest) *models.Faucet {
	return &models.Faucet{
		UUID:            uuid.NewV4().String(),
		Name:            request.Name,
		ChainRule:       request.ChainRule,
		CreditorAccount: request.CreditorAccount,
		MaxBalance:      request.MaxBalance,
		Amount:          request.Amount,
		Cooldown:        request.Cooldown,
	}
}
