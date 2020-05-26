package controllers

import (
	"encoding/json"
	"net/http"

	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/formatters"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions"

	"github.com/gorilla/mux"
)

type TransactionsController struct {
	ucs transactions.UseCases
}

func NewTransactionsController(useCases transactions.UseCases) *TransactionsController {
	return &TransactionsController{
		ucs: useCases,
	}
}

// Add routes to router
func (c *TransactionsController) Append(router *mux.Router) {
	router.Methods(http.MethodPost).Path("/transactions/{chainUUID}/send").HandlerFunc(c.Send)
}

// @Summary Creates and sends a new contract transaction
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 400
// @Failure 409
// @Failure 422
// @Failure 500
// @Router /transactions/{chainUUID}/send [post]
func (c *TransactionsController) Send(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	txRequest := &types.SendTransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, txRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	chainUUID := mux.Vars(request)["chainUUID"]
	tenantID := multitenancy.TenantIDFromContext(ctx)
	txReq := formatters.FormatSendTxRequest(txRequest, chainUUID)
	txResponse, err := c.ucs.SendTransaction().Execute(ctx, txReq, tenantID)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	response, err := formatters.FormatTxResponse(txResponse)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(response)
}
