package faucets

import (
	"encoding/json"
	"net/http"

	"github.com/asaskevich/govalidator"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/chainregistry"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/store/models"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
)

var _ models.Faucet
var _ chainregistry.Faucet

// @Summary Retrieves a list of all registered faucet
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200 {array} models.Faucet
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /faucets [get]
func (h *controller) GetFaucets(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	faucets, err := h.getFaucetsUC.Execute(
		request.Context(),
		multitenancy.AllowedTenantsFromContext(request.Context()),
		httputil.ToFilters(request.URL.Query()),
	)

	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	if len(faucets) == 0 {
		faucets = []*models.Faucet{}
	}

	_ = json.NewEncoder(rw).Encode(faucets)
}

// @Summary Retrieves a faucet by ID
// @Produce json
// @Param uuid path string true "ID of the faucet"
// @Success 200 {object} models.Faucet
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Faucet not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /faucets/{uuid} [get]
func (h *controller) GetFaucet(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	faucet, err := h.getFaucetUC.Execute(
		request.Context(),
		mux.Vars(request)["uuid"],
		multitenancy.AllowedTenantsFromContext(request.Context()),
	)

	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(faucet)
}

// @Summary Retrieve faucet candidate for provided sender and chain
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param chain_uuid query string true "chain uuid to calculate faucet candidate"
// @Param account query string true "hex address of account sender"
// @Success 200 {object} chainregistry.Faucet
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /faucets/candidate [get]
func (h *controller) GetFaucetCandidate(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chainUUID := req.URL.Query().Get("chain_uuid")
	if !govalidator.IsUUID(chainUUID) {
		err := errors.DataError("invalid \"chain_uuid\" value. Expected uuid, found %s", chainUUID)
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	account := req.URL.Query().Get("account")
	if !ethcommon.IsHexAddress(account) {
		err := errors.DataError("invalid \"account\" value. Expected hex, found %s", account)
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	tenants := multitenancy.AllowedTenantsFromContext(req.Context())
	faucet, err := h.faucetCandidateUC.Execute(req.Context(), ethcommon.HexToAddress(account), chainUUID, tenants)
	if err != nil {
		if errors.IsFaucetWarning(err) {
			err = errors.NotFoundError(err.Error())
		}

		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(faucet)
}
