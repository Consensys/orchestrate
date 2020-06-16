package controllers

import (
	"encoding/json"
	"net/http"

	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/formatters"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/chains"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions"

	"github.com/gorilla/mux"
)

type TransactionsController struct {
	txUcs    transactions.UseCases
	chainUcs chains.UseCases
}

func NewTransactionsController(txUcs transactions.UseCases, chainUcs chains.UseCases) *TransactionsController {
	return &TransactionsController{
		txUcs:    txUcs,
		chainUcs: chainUcs,
	}
}

// Add routes to router
func (c *TransactionsController) Append(router *mux.Router) {
	router.Methods(http.MethodPost).Path("/transactions/send").HandlerFunc(c.Send)
	router.Methods(http.MethodPost).Path("/transactions/send-raw").HandlerFunc(c.SendRaw)
	router.Methods(http.MethodPost).Path("/transactions/send-transfer").HandlerFunc(c.SendTransfer)
	router.Methods(http.MethodPost).Path("/transactions/deploy-contract").HandlerFunc(c.DeployContract)
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
// @Router /transactions/send [post]
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

	tenantID := multitenancy.TenantIDFromContext(ctx)
	txReq := formatters.FormatSendTxRequest(txRequest)

	chain, err := c.chainUcs.GetChainByName().Execute(ctx, txRequest.ChainName, tenantID)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	txResponse, err := c.txUcs.SendContractTransaction().Execute(ctx, txReq, chain.UUID, tenantID)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txResponse, chain.Name))
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
// @Router /transactions/deploy-contract [post]
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

	tenantID := multitenancy.TenantIDFromContext(ctx)
	txReq := formatters.FormatDeployContractRequest(txRequest)

	chain, err := c.chainUcs.GetChainByName().Execute(ctx, txRequest.ChainName, tenantID)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	txResponse, err := c.txUcs.SendDeployTransaction().Execute(ctx, txReq, chain.UUID, tenantID)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txResponse, chain.Name))
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
// @Router /transactions/send-raw [post]
func (c *TransactionsController) SendRaw(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	txRequest := &types.RawTransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, txRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	tenantID := multitenancy.TenantIDFromContext(ctx)
	txReq := formatters.FormatSendRawRequest(txRequest)

	chain, err := c.chainUcs.GetChainByName().Execute(ctx, txRequest.ChainName, tenantID)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	txResponse, err := c.txUcs.SendTransaction().Execute(ctx, txReq, "", chain.UUID, tenantID)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txResponse, chain.Name))
}

// @Summary Creates and sends a transfer transaction
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 400
// @Failure 409
// @Failure 422
// @Failure 500
// @Router /transactions/send-transfer [post]
func (c *TransactionsController) SendTransfer(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	txRequest := &types.TransferRequest{}
	err := jsonutils.UnmarshalBody(request.Body, txRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	tenantID := multitenancy.TenantIDFromContext(ctx)
	txReq := formatters.FormatSendTransferRequest(txRequest)

	chain, err := c.chainUcs.GetChainByName().Execute(ctx, txRequest.ChainName, tenantID)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	txResponse, err := c.txUcs.SendTransaction().Execute(ctx, txReq, "", chain.UUID, tenantID)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txResponse, chain.Name))
}
