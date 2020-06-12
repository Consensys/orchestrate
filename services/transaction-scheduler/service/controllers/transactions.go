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
	router.Methods(http.MethodPost).Path("/transactions/{chainUUID}/send-raw").HandlerFunc(c.SendRaw)
	router.Methods(http.MethodPost).Path("/transactions/{chainUUID}/deploy-contract").HandlerFunc(c.DeployContract)
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
	if err := jsonutils.UnmarshalBody(request.Body, txRequest); err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := txRequest.Params.PrivateTransactionParams.Validate(); err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	chainUUID := mux.Vars(request)["chainUUID"]
	tenantID := multitenancy.TenantIDFromContext(ctx)
	txReq := formatters.FormatSendTxRequest(txRequest)

	txResponse, err := c.ucs.SendContractTransaction().Execute(ctx, txReq, chainUUID, tenantID)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txResponse))
}

// @Summary Creates and sends a new contract deployment transaction
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 400
// @Failure 409
// @Failure 422
// @Failure 500
// @Router /transactions/{chainUUID}/deploy-contract [post]
func (c *TransactionsController) DeployContract(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	txRequest := &types.DeployContractRequest{}
	if err := jsonutils.UnmarshalBody(request.Body, txRequest); err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := txRequest.Params.PrivateTransactionParams.Validate(); err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	chainUUID := mux.Vars(request)["chainUUID"]
	tenantID := multitenancy.TenantIDFromContext(ctx)
	txReq := formatters.FormatDeployContractRequest(txRequest)

	txResponse, err := c.ucs.SendDeployTransaction().Execute(ctx, txReq, chainUUID, tenantID)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txResponse))
}

// @Summary Creates and sends a raw transaction
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 400
// @Failure 409
// @Failure 422
// @Failure 500
// @Router /transactions/{chainUUID}/send-raw [post]
func (c *TransactionsController) SendRaw(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	txRequest := &types.RawTransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, txRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	chainUUID := mux.Vars(request)["chainUUID"]
	tenantID := multitenancy.TenantIDFromContext(ctx)
	txReq := formatters.FormatSendRawRequest(txRequest)
	txResponse, err := c.ucs.SendTransaction().Execute(ctx, txReq, "", chainUUID, tenantID)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txResponse))
}
