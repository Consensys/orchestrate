package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/consensys/orchestrate/pkg/errors"

	"github.com/consensys/orchestrate/pkg/types/api"

	jsonutils "github.com/consensys/orchestrate/pkg/encoding/json"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/types/formatters"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/gorilla/mux"
)

var _ entities.ETHTransactionParams

const (
	IdempotencyKeyHeader = "X-Idempotency-Key"
	MinBoostValue        = 0.1
)

type TransactionsController struct {
	ucs usecases.TransactionUseCases
}

func NewTransactionsController(ucs usecases.TransactionUseCases) *TransactionsController {
	return &TransactionsController{
		ucs: ucs,
	}
}

// Add routes to router
func (c *TransactionsController) Append(router *mux.Router) {
	router.Methods(http.MethodPost).Path("/transactions/send").
		Handler(http.HandlerFunc(c.send))
	router.Methods(http.MethodPost).Path("/transactions/send-raw").
		Handler(http.HandlerFunc(c.sendRaw))
	router.Methods(http.MethodPost).Path("/transactions/transfer").
		Handler(http.HandlerFunc(c.transfer))
	router.Methods(http.MethodPost).Path("/transactions/deploy-contract").
		Handler(http.HandlerFunc(c.deployContract))
	router.Methods(http.MethodGet).Path("/transactions/{uuid}").
		Handler(http.HandlerFunc(c.getOne))
	router.Methods(http.MethodPut).Path("/transactions/{uuid}/speed-up").
		Handler(http.HandlerFunc(c.speedUp))
	router.Methods(http.MethodPut).Path("/transactions/{uuid}/call-off").
		Handler(http.HandlerFunc(c.callOff))
	router.Methods(http.MethodGet).Path("/transactions").
		Handler(http.HandlerFunc(c.search))
}

// @Summary Creates and sends a new contract transaction
// @Description Creates and executes a new smart contract transaction request
// @Description The transaction can be private (Tessera, EEA).
// @Description The transaction can be a One Time Key transaction in 0 gas private networks
// @Tags Transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body api.SendTransactionRequest{params=api.TransactionParams{gasPricePolicy=api.GasPriceParams{retryPolicy=api.RetryParams}}} true "Contract transaction request"
// @Success 202 {object} api.TransactionResponse "Created contract transaction request"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 409 {object} httputil.ErrorResponse "Already existing transaction"
// @Failure 422 {object} httputil.ErrorResponse "Unprocessable parameters were sent"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /transactions/send [post]
func (c *TransactionsController) send(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	txRequest := &api.SendTransactionRequest{}
	if err := jsonutils.UnmarshalBody(request.Body, txRequest); err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := txRequest.Params.Validate(); err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	txReq := formatters.FormatSendTxRequest(txRequest, request.Header.Get(IdempotencyKeyHeader))
	txResponse, err := c.ucs.SendContractTransaction().Execute(ctx, txReq, multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txResponse))
}

// @Summary Creates and sends a new contract deployment
// @Description Creates and executes a new contract deployment request
// @Description The transaction can be private (Tessera, EEA).
// @Description The transaction can be a One Time Key transaction in 0 gas private networks
// @Tags Transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body api.DeployContractRequest{params=api.DeployContractParams{gasPricePolicy=api.GasPriceParams{retryPolicy=api.RetryParams}}} true "Deployment transaction request"
// @Success 202 {object} api.TransactionResponse "Created deployment transaction request"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 409 {object} httputil.ErrorResponse "Already existing transaction"
// @Failure 422 {object} httputil.ErrorResponse "Unprocessable parameters were sent"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /transactions/deploy-contract [post]
func (c *TransactionsController) deployContract(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	txRequest := &api.DeployContractRequest{}
	if err := jsonutils.UnmarshalBody(request.Body, txRequest); err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err := txRequest.Params.Validate(); err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	txReq := formatters.FormatDeployContractRequest(txRequest, request.Header.Get(IdempotencyKeyHeader))
	txResponse, err := c.ucs.SendDeployTransaction().Execute(ctx, txReq, multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txResponse))
}

// @Summary Creates and sends a raw transaction
// @Description Creates and executes a new raw transaction request
// @Tags Transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body api.RawTransactionRequest{params=api.RawTransactionParams{retryPolicy=api.IntervalRetryParams}} true "Raw transaction request"
// @Success 202 {object} api.TransactionResponse "Created raw transaction request"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 409 {object} httputil.ErrorResponse "Already existing transaction"
// @Failure 422 {object} httputil.ErrorResponse "Unprocessable parameters were sent"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /transactions/send-raw [post]
func (c *TransactionsController) sendRaw(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	txRequest := &api.RawTransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, txRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	txReq := formatters.FormatSendRawRequest(txRequest, request.Header.Get(IdempotencyKeyHeader))
	txResponse, err := c.ucs.SendTransaction().Execute(ctx, txReq, nil, multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txResponse))
}

// @Summary Creates and sends a transfer transaction
// @Description Creates and executes a new transfer request
// @Tags Transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body api.TransferRequest{params=api.TransferParams{gasPricePolicy=api.GasPriceParams{retryPolicy=api.RetryParams}}} true "Transfer transaction request"
// @Success 202 {object} api.TransactionResponse "Created transfer transaction request"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 409 {object} httputil.ErrorResponse "Already existing transaction"
// @Failure 422 {object} httputil.ErrorResponse "Unprocessable parameters were sent"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /transactions/transfer [post]
func (c *TransactionsController) transfer(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	txRequest := &api.TransferRequest{}
	err := jsonutils.UnmarshalBody(request.Body, txRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err = txRequest.Params.Validate(); err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	txReq := formatters.FormatTransferRequest(txRequest, request.Header.Get(IdempotencyKeyHeader))
	txResponse, err := c.ucs.SendTransaction().Execute(ctx, txReq, nil, multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txResponse))
}

// @Summary Fetch a transaction request by uuid
// @Description Fetch a single transaction request by uuid
// @Tags Transactions
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "UUID of the transaction request"
// @Success 200 {object} api.TransactionResponse "Transaction request found"
// @Failure 404 {object} httputil.ErrorResponse "Transaction request not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /transactions/{uuid} [get]
func (c *TransactionsController) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	uuid := mux.Vars(request)["uuid"]

	txRequest, err := c.ucs.GetTransaction().Execute(ctx, uuid, multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(txRequest))
}

// @Summary      Search transaction requests by provided filters
// @Description  Get a list of filtered transaction requests
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Security     JWTAuth
// @Param        idempotency_keys  query     []string                 false  "List of idempotency keys"  collectionFormat(csv)
// @Param		 limit	  query     int					 false  "Maximum response size"
// @Param		 page	  query     int					 false  "Page within the entire responses set"
// @Success      200               {array}   api.TransactionResponse  "List of transaction requests found"
// @Failure      400               {object}  infra.ErrorResponse      "Invalid filter in the request"
// @Failure      500               {object}  infra.ErrorResponse      "Internal server error"
// @Router       /transactions [get]
func (c *TransactionsController) search(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	filters, err := formatters.FormatTransactionsFilterRequest(request)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	// increase Limit to test more items available
	if filters.Pagination.Limit > 0 {
		filters.Pagination.Limit++
	} else {
		filters.Pagination.Limit = api.DefaultTransactionPageSize + 1
	}

	txs, err := c.ucs.SearchTransactions().Execute(ctx, filters, multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	// signal when more items are available
	hasMore := len(txs) == filters.Pagination.Limit

	response := &api.TransactionSearchResponse{
		HasMore: hasMore}

	if response.HasMore {
		txs = txs[:filters.Pagination.Limit-1]
	}

	var resTxs []*api.TransactionResponse
	for _, tx := range txs {
		resTxs = append(resTxs, formatters.FormatTxResponse(tx))
	}

	response.Transactions = resTxs

	_ = json.NewEncoder(rw).Encode(response)
}

// @Summary      Speed up transaction timeIncrease transaction gas price
// @Description  Speed up transaction time by an increase its gas price
// @Tags         Transactions
// @Produce      json
// @Security     ApiKeyAuth
// @Security     JWTAuth
// @Param        uuid  path      string               true  "UUID of the transaction request"
// @Param        boost  query     float64              false  "gas price increment percentage, default value is 0.05 (5%)"
// @Failure      404   {object}  infra.ErrorResponse  "Transaction request not found"
// @Failure      500   {object}  infra.ErrorResponse  "Internal server error"
// @Router       /transactions/{uuid}/speed-up [put]
func (c *TransactionsController) speedUp(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := req.Context()
	var boostValue = MinBoostValue // Default
	var err error

	uuid := mux.Vars(req)["uuid"]
	boostInput := req.URL.Query().Get("boost")
	if boostInput != "" {
		boostValue, err = strconv.ParseFloat(boostInput, 64)
		if err != nil {
			httputil.WriteHTTPErrorResponse(rw, errors.InvalidFormatError("expected float as incrementing value"))
			return
		}
		if boostValue < MinBoostValue {
			err = errors.InvalidFormatError("minimum boost value is %f", MinBoostValue)
			httputil.WriteHTTPErrorResponse(rw, err)
			return
		}
	}

	tx, err := c.ucs.SpeedUp().Execute(ctx, uuid, boostValue, multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(tx))
}

// @Summary      Call off transaction
// @Description  Call off transaction by sending empty data transaction with same nonce and additional 10% more gas
// @Tags         Transactions
// @Produce      json
// @Security     ApiKeyAuth
// @Security     JWTAuth
// @Param        uuid  path      string               true  "UUID of the transaction request"
// @Failure      404   {object}  infra.ErrorResponse  "Transaction request not found"
// @Failure      500   {object}  infra.ErrorResponse  "Internal server error"
// @Router       /transactions/{uuid}/call-off [put]
func (c *TransactionsController) callOff(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	uuid := mux.Vars(request)["uuid"]

	tx, err := c.ucs.CallOff().Execute(ctx, uuid, multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatTxResponse(tx))
}
