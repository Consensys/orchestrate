package controllers

import (
	"encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions"
	"net/http"

	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"

	"github.com/gorilla/mux"
)

type TransactionsController struct {
	sendTxUseCase transactions.SendTxUseCase
}

func NewTransactionsController(sendTxUseCase transactions.SendTxUseCase) *TransactionsController {
	return &TransactionsController{
		sendTxUseCase: sendTxUseCase,
	}
}

// Add routes to router
func (c *TransactionsController) Append(router *mux.Router) {
	router.Methods(http.MethodPost).Path("/transactions/send").HandlerFunc(c.Send)
}

// @Summary Creates and sends a new transaction
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200 types.TransactionRequest
// @Failure 400 error
// @Failure 409 error
// @Failure 422 error
// @Failure 500 error
// @Router /transactions/send [post]
func (c *TransactionsController) Send(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	txRequest := &types.TransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, txRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	tenantID := multitenancy.TenantIDFromContext(ctx)
	transactionResponse, err := c.sendTxUseCase.Execute(ctx, txRequest, tenantID)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(transactionResponse)
}
