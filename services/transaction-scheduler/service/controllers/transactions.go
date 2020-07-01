package controllers

import (
	"encoding/json"
	"net/http"

	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/formatters"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/service/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/chains"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/transactions"

	"github.com/gorilla/mux"
)

const (
	IdempotencyKeyHeader = "X-Idempotency-Key"
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
	router.Methods(http.MethodPost).Path("/transactions/send").
		Handler(idempotencyKeyMiddleware(http.HandlerFunc(c.send)))
	router.Methods(http.MethodPost).Path("/transactions/send-raw").
		Handler(idempotencyKeyMiddleware(http.HandlerFunc(c.sendRaw)))
	router.Methods(http.MethodPost).Path("/transactions/transfer").
		Handler(idempotencyKeyMiddleware(http.HandlerFunc(c.transfer)))
	router.Methods(http.MethodPost).Path("/transactions/deploy-contract").
		Handler(idempotencyKeyMiddleware(http.HandlerFunc(c.deployContract)))
	router.Methods(http.MethodGet).Path("/transactions/{uuid}").
		Handler(http.HandlerFunc(c.getOne))
	router.Methods(http.MethodGet).Path("/transactions").
		Handler(http.HandlerFunc(c.search))
}

func idempotencyKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(IdempotencyKeyHeader) == "" {
			r.Header.Set(IdempotencyKeyHeader, utils.RandomString(16))
		}

		next.ServeHTTP(w, r)
	})
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
func (c *TransactionsController) send(rw http.ResponseWriter, request *http.Request) {
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
	txReq := formatters.FormatSendTxRequest(txRequest, request.Header.Get(IdempotencyKeyHeader))

	chain, err := c.chainUcs.GetChainByName().Execute(ctx, txRequest.ChainName,
		multitenancy.AllowedTenantsFromContext(ctx))
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
// @Router /transactions/deploy-contract [post]
func (c *TransactionsController) deployContract(rw http.ResponseWriter, request *http.Request) {
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

	txReq := formatters.FormatDeployContractRequest(txRequest, request.Header.Get(IdempotencyKeyHeader))

	chain, err := c.chainUcs.GetChainByName().Execute(ctx, txRequest.ChainName,
		multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	txResponse, err := c.txUcs.SendDeployTransaction().Execute(ctx, txReq, chain.UUID,
		multitenancy.TenantIDFromContext(ctx))
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
// @Router /transactions/send-raw [post]
func (c *TransactionsController) sendRaw(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	txRequest := &types.RawTransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, txRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	txReq := formatters.FormatSendRawRequest(txRequest, request.Header.Get(IdempotencyKeyHeader))

	chain, err := c.chainUcs.GetChainByName().Execute(ctx, txRequest.ChainName,
		multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	txResponse, err := c.txUcs.SendTransaction().Execute(ctx, txReq, "", chain.UUID,
		multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txResponse))
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
// @Router /transactions/transfer [post]
func (c *TransactionsController) transfer(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	txRequest := &types.TransferRequest{}
	err := jsonutils.UnmarshalBody(request.Body, txRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	txReq := formatters.FormatSendTransferRequest(txRequest, request.Header.Get(IdempotencyKeyHeader))

	chain, err := c.chainUcs.GetChainByName().Execute(ctx, txRequest.ChainName,
		multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	txResponse, err := c.txUcs.SendTransaction().Execute(ctx, txReq, "", chain.UUID,
		multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txResponse))
}

// @Summary Gets a transaction request by uuid
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 404
// @Failure 500
// @Router /transactions/{uuid} [get]
func (c *TransactionsController) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	uuid := mux.Vars(request)["uuid"]

	txRequest, err := c.txUcs.GetTransaction().Execute(ctx, uuid, multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txRequest))
}

// @Summary Searches transaction requests by filter
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200
// @Failure 400
// @Failure 500
// @Router /transactions [get]
func (c *TransactionsController) search(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	filters, err := formatters.FormatTransactionsFilterRequest(request)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	txRequests, err := c.txUcs.SearchTransactions().Execute(ctx, filters, multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	var responses []*types.TransactionResponse
	for _, txRequest := range txRequests {
		responses = append(responses, formatters.FormatTxResponse(txRequest))
	}

	_ = json.NewEncoder(rw).Encode(responses)
}
