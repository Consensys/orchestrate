package faucets

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/multitenancy"
)

type PostRequest struct {
	Name            string `json:"name" validate:"required"`
	ChainRule       string `json:"chainRule,omitempty" validate:"required"`
	CreditorAccount string `json:"creditorAccount,omitempty" validate:"required,eth_addr"`
	MaxBalance      string `json:"maxBalance,omitempty" validate:"required,isBig"`
	Amount          string `json:"amount,omitempty" validate:"required,isBig"`
	Cooldown        string `json:"cooldown,omitempty" validate:"required,isDuration"`
}

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
func (h Handler) postFaucet(rw http.ResponseWriter, request *http.Request) {
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

	faucet := &types.Faucet{
		UUID:            uuid.NewV4().String(),
		Name:            faucetRequest.Name,
		TenantID:        tenantID,
		ChainRule:       faucetRequest.ChainRule,
		CreditorAccount: faucetRequest.CreditorAccount,
		MaxBalance:      faucetRequest.MaxBalance,
		Amount:          faucetRequest.Amount,
		Cooldown:        faucetRequest.Cooldown,
	}

	err = h.store.RegisterFaucet(request.Context(), faucet)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(faucet)
}
