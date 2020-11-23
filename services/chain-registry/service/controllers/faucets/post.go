package faucets

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"
)

// @Summary Registers a new faucet
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body PostRequest true "Faucet registration request"
// @Success 200 {object} models.Faucet
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /faucets [post]
func (h *controller) PostFaucet(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	faucetRequest := &PostRequest{}

	err := jsonutils.UnmarshalBody(request.Body, faucetRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	faucet := parsePostRequestToFaucet(faucetRequest)
	faucet.TenantID = multitenancy.TenantIDFromContext(request.Context())

	err = h.registerFaucetUC.Execute(request.Context(), faucet)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(rw).Encode(faucet)
}

func parsePostRequestToFaucet(request *PostRequest) *models.Faucet {
	return &models.Faucet{
		UUID:            uuid.Must(uuid.NewV4()).String(),
		Name:            request.Name,
		ChainRule:       request.ChainRule,
		CreditorAccount: request.CreditorAccount,
		MaxBalance:      request.MaxBalance,
		Amount:          request.Amount,
		Cooldown:        request.Cooldown,
	}
}
