package faucets

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
)

type PatchRequest struct {
	Name            string `json:"name,omitempty" validate:"omitempty"`
	ChainRule       string `json:"chainRule,omitempty" validate:"omitempty"`
	CreditorAccount string `json:"creditorAccount,omitempty" validate:"omitempty,eth_addr"`
	MaxBalance      string `json:"maxBalance,omitempty" validate:"omitempty,isBig"`
	Amount          string `json:"amount,omitempty" validate:"omitempty,isBig"`
	Cooldown        string `json:"cooldown,omitempty" validate:"omitempty,isDuration"`
}

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
func (h Handler) patchFaucetByUUID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	faucetRequest := &PatchRequest{}

	err := utils.UnmarshalBody(request.Body, faucetRequest)
	if err != nil {
		utils.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	faucet := &types.Faucet{
		UUID:            mux.Vars(request)["uuid"],
		Name:            faucetRequest.Name,
		ChainRule:       faucetRequest.ChainRule,
		CreditorAccount: faucetRequest.CreditorAccount,
		MaxBalance:      faucetRequest.MaxBalance,
		Amount:          faucetRequest.Amount,
		Cooldown:        faucetRequest.Cooldown,
	}

	tenantID := multitenancy.TenantIDFromContext(request.Context())
	if tenantID != "" {
		faucet.TenantID = tenantID
	}

	err = h.store.UpdateFaucetByUUID(request.Context(), faucet)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(faucet)
}
