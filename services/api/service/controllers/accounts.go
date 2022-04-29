package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	jsonutils "github.com/consensys/orchestrate/pkg/encoding/json"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/formatters"
	"github.com/consensys/orchestrate/pkg/utils"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/quorum-key-manager/pkg/client"
	qkmstoretypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
	qkmutilstypes "github.com/consensys/quorum-key-manager/src/utils/api/types"
	"github.com/gorilla/mux"
)

type AccountsController struct {
	ucs              usecases.AccountUseCases
	keyManagerClient client.KeyManagerClient
	storeName        string
}

func NewAccountsController(accountUCs usecases.AccountUseCases, keyManagerClient client.KeyManagerClient, qkmStoreID string) *AccountsController {
	return &AccountsController{
		accountUCs,
		keyManagerClient,
		qkmStoreID,
	}
}

// Append Add routes to router
func (c *AccountsController) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/accounts").HandlerFunc(c.search)
	router.Methods(http.MethodPost).Path("/accounts").HandlerFunc(c.create)
	router.Methods(http.MethodPost).Path("/accounts/import").HandlerFunc(c.importKey)
	router.Methods(http.MethodGet).Path("/accounts/{address}").HandlerFunc(c.getOne)
	router.Methods(http.MethodDelete).Path("/accounts/{address}").HandlerFunc(c.deleteOne)
	router.Methods(http.MethodPatch).Path("/accounts/{address}").HandlerFunc(c.update)
	router.Methods(http.MethodPost).Path("/accounts/{address}/sign-message").HandlerFunc(c.signMessage)
	router.Methods(http.MethodPost).Path("/accounts/{address}/sign-typed-data").HandlerFunc(c.signTypedData)
	router.Methods(http.MethodPost).Path("/accounts/verify-message").HandlerFunc(c.verifyMessageSignature)
	router.Methods(http.MethodPost).Path("/accounts/verify-typed-data").HandlerFunc(c.verifyTypedDataSignature)
}

// @Summary Creates a new Account
// @Description Creates a new Account
// @Tags Accounts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body api.CreateAccountRequest true "Account creation request"
// @Success 200 {object} api.AccountResponse "Account object"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts [post]
func (c *AccountsController) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	req := &api.CreateAccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, req)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	acc, err := c.ucs.CreateAccount().Execute(ctx, formatters.FormatCreateAccountRequest(req, c.storeName), nil, req.Chain,
		multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(acc))
}

// @Summary Fetch an account by address
// @Description Fetch a single account by address
// @Tags Accounts
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param address path string true "selected account address"
// @Success 200 {object} api.AccountResponse "Account found"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/{address} [get]
func (c *AccountsController) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	address, err := utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	acc, err := c.ucs.GetAccount().Execute(ctx, *address, multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(acc))
}

// @Summary Delete an account by address
// @Description Delete a single account by address
// @Tags Accounts
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param address path string true "deleted account address"
// @Success 204
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/{address} [delete]
func (c *AccountsController) deleteOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	address, err := utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.ucs.DeleteAccount().Execute(ctx, *address, multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Search accounts by provided filters
// @Description Get a list of filtered accounts
// @Tags Accounts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param aliases query []string false "List of account aliases" collectionFormat(csv)
// @Success 200 {array} api.AccountResponse "List of identities found"
// @Failure 400 {object} httputil.ErrorResponse "Invalid filter in the request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts [get]
func (c *AccountsController) search(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	filters, err := formatters.FormatAccountFilterRequest(request)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	accs, err := c.ucs.SearchAccounts().Execute(ctx, filters, multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	response := []*api.AccountResponse{}
	for _, iden := range accs {
		response = append(response, formatters.FormatAccountResponse(iden))
	}

	_ = json.NewEncoder(rw).Encode(response)
}

// @Summary Creates a new Account by importing a private key
// @Description Creates a new Account by importing a private key
// @Tags Accounts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body api.ImportAccountRequest true "Account creation request"
// @Success 200 {object} api.AccountResponse "Account object"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 422 {object} httputil.ErrorResponse "Unprocessable entity"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 405 {object} httputil.ErrorResponse "Not allowed"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/import [post]
func (c *AccountsController) importKey(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	req := &api.ImportAccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, req)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	acc, err := c.ucs.CreateAccount().Execute(ctx, formatters.FormatImportAccountRequest(req, c.storeName), req.PrivateKey, req.Chain,
		multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(acc))
}

// @Summary Update account by Address
// @Description Update a specific account by Address
// @Tags Accounts
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body api.UpdateAccountRequest true "Account update request"
// @Param address path string true "selected account address"
// @Success 200 {object} api.AccountResponse "Account found"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/{address} [patch]
func (c *AccountsController) update(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	accRequest := &api.UpdateAccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, accRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	acc := formatters.FormatUpdateAccountRequest(accRequest)
	address, err := utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	acc.Address = *address

	accRes, err := c.ucs.UpdateAccount().Execute(ctx, acc, multitenancy.UserInfoValue(ctx))

	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(accRes))
}

// @Summary Sign Message (EIP-191)
// @Description Sign message, following EIP-191, data using selected account
// @Tags Accounts
// @Accept json
// @Produce text/plain
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body api.SignMessageRequest true "Payload to sign"
// @Param address path string true "selected account address"
// @Success 200 {string} string "Signed payload"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/{address}/sign-message [post]
func (c *AccountsController) signMessage(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	payloadRequest := &api.SignMessageRequest{}
	err := jsonutils.UnmarshalBody(request.Body, payloadRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	address, err := utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.ucs.GetAccount().Execute(ctx, *address, multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteError(rw, fmt.Sprintf("account %s was not found", address), http.StatusBadRequest)
		return
	}

	qkmStoreID := payloadRequest.StoreID
	if qkmStoreID == "" {
		qkmStoreID = c.storeName
	}

	signature, err := c.keyManagerClient.SignMessage(request.Context(), qkmStoreID, address.Hex(), &qkmstoretypes.SignMessageRequest{
		Message: payloadRequest.Message,
	})
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(signature))
}

// @Summary Signs typed data using an existing account following the EIP-712 standard
// @Description Signs typed data using ECDSA and the private key of an existing account following the EIP-712 standard
// @Tags Accounts
// @Accept json
// @Produce text/plain
// @Param request body api.SignTypedDataRequest true "Typed data to sign"
// @Param address path string true "selected account address"
// @Success 200 {string} string "Signed payload"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 422 {object} httputil.ErrorResponse "Invalid parameters"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/{address}/sign-typed-data [post]
func (c *AccountsController) signTypedData(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	signRequest := &api.SignTypedDataRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	address, err := utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.ucs.GetAccount().Execute(ctx, *address, multitenancy.UserInfoValue(ctx))
	if err != nil {
		httputil.WriteError(rw, fmt.Sprintf("account %s was not found", address), http.StatusBadRequest)
		return
	}

	qkmStoreID := signRequest.StoreID
	if qkmStoreID == "" {
		qkmStoreID = c.storeName
	}

	signature, err := c.keyManagerClient.SignTypedData(ctx, qkmStoreID, address.Hex(), &qkmstoretypes.SignTypedDataRequest{
		DomainSeparator: signRequest.DomainSeparator,
		Types:           signRequest.Types,
		Message:         signRequest.Message,
		MessageType:     signRequest.MessageType,
	})
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(signature))
}

// @Summary Verifies the signature of a typed data message following the EIP-712 standard
// @Description Verifies if a typed data message has been signed by the Ethereum account passed as argument following the EIP-712 standard
// @Tags Accounts
// @Accept json
// @Param request body qkmutilstypes.VerifyTypedDataRequest true "Typed data to sign"
// @Success 204
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 422 {object} httputil.ErrorResponse "Invalid parameters"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/verify-typed-data [post]
func (c *AccountsController) verifyTypedDataSignature(rw http.ResponseWriter, request *http.Request) {
	verifyRequest := &qkmutilstypes.VerifyTypedDataRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.keyManagerClient.VerifyTypedData(request.Context(), verifyRequest)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Verifies the signature of a message (EIP-191)
// @Description Verifies if a message has been signed by the Ethereum account passed as argument
// @Tags Accounts
// @Accept json
// @Param request body qkmutilstypes.VerifyRequest true "signature and message to verify"
// @Success 204
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 422 {object} httputil.ErrorResponse "Invalid parameters"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/verify-message [post]
func (c *AccountsController) verifyMessageSignature(rw http.ResponseWriter, request *http.Request) {
	verifyRequest := &qkmutilstypes.VerifyRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.keyManagerClient.VerifyMessage(request.Context(), verifyRequest)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
