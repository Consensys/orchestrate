package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	qkm "github.com/consensys/orchestrate/pkg/quorum-key-manager"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/services/api/service/formatters"
	"github.com/consensys/quorum-key-manager/pkg/client"
	qkmtypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/ethereum/go-ethereum/common/hexutil"

	jsonutils "github.com/consensys/orchestrate/pkg/encoding/json"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/utils"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/gorilla/mux"
)

type AccountsController struct {
	ucs              usecases.AccountUseCases
	keyManagerClient client.KeyManagerClient
	storeName        string
}

func NewAccountsController(accountUCs usecases.AccountUseCases, keyManagerClient client.KeyManagerClient) *AccountsController {
	return &AccountsController{
		accountUCs,
		keyManagerClient,
		qkm.GlobalStoreName(),
	}
}

// Add routes to router
func (c *AccountsController) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/accounts").HandlerFunc(c.search)
	router.Methods(http.MethodPost).Path("/accounts").HandlerFunc(c.create)
	router.Methods(http.MethodPost).Path("/accounts/import").HandlerFunc(c.importKey)
	router.Methods(http.MethodGet).Path("/accounts/{address}").HandlerFunc(c.getOne)
	router.Methods(http.MethodPatch).Path("/accounts/{address}").HandlerFunc(c.update)
	router.Methods(http.MethodPost).Path("/accounts/{address}/sign-message").HandlerFunc(c.signMessage)
	router.Methods(http.MethodPost).Path("/accounts/{address}/sign-typed-data").HandlerFunc(c.signTypedData)
	router.Methods(http.MethodPost).Path("/accounts/verify-message").HandlerFunc(c.verifyMessageSignature)
	router.Methods(http.MethodPost).Path("/accounts/verify-typed-data").HandlerFunc(c.verifyTypedDataSignature)
}

// @Summary Creates a new Account
// @Description Creates a new Account
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

	acc := formatters.FormatCreateAccountRequest(req)
	acc, err = c.ucs.CreateAccount().Execute(ctx, acc, nil, req.Chain, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(acc))
}

// @Summary Fetch a account by address
// @Description Fetch a single account by address
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

	acc, err := c.ucs.GetAccount().Execute(ctx, address, multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(acc))
}

// @Summary Search accounts by provided filters
// @Description Get a list of filtered accounts
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

	accs, err := c.ucs.SearchAccounts().Execute(ctx, filters, multitenancy.AllowedTenantsFromContext(ctx))
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

	var bPrivKey []byte
	if req.PrivateKey != "" {
		bPrivKey, err = hexutil.Decode("0x" + req.PrivateKey)
		if err != nil {
			httputil.WriteError(rw, "invalid private key format", http.StatusBadRequest)
			return
		}
	}

	accResp := formatters.FormatImportAccountRequest(req)
	accResp, err = c.ucs.CreateAccount().Execute(ctx, accResp, bPrivKey, req.Chain, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(accResp))
}

// @Summary Update account by Address
// @Description Update a specific account by Address
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
	acc.Address, err = utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	accRes, err := c.ucs.UpdateAccount().Execute(ctx, acc, multitenancy.AllowedTenantsFromContext(ctx))

	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(accRes))
}

// @Summary Sign Message (EIP-191)
// @Description Sign message, following EIP-191, data using selected account
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

	tenants := utils.AllowedTenants(multitenancy.TenantIDFromContext(ctx))
	_, err = c.ucs.GetAccount().Execute(ctx, address, tenants)
	if err != nil {
		httputil.WriteError(rw, fmt.Sprintf("account %s was not found", address), http.StatusBadRequest)
		return
	}

	signature, err := c.keyManagerClient.SignMessage(request.Context(), c.storeName, address, &qkmtypes.SignMessageRequest{
		Message: hexutil.MustDecode(payloadRequest.Message),
	})
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(signature))
}

// @Summary Signs typed data using an existing account following the EIP-712 standard
// @Description Signs typed data using ECDSA and the private key of an existing account following the EIP-712 standard
// @Accept json
// @Produce text/plain
// @Param request body api.SignTypedDataRequest{domainSeparator=qkmtypes.DomainSeparator,types=map[string]qkmtypes.Type} true "Typed data to sign"
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

	tenants := utils.AllowedTenants(multitenancy.TenantIDFromContext(ctx))
	_, err = c.ucs.GetAccount().Execute(ctx, address, tenants)
	if err != nil {
		httputil.WriteError(rw, fmt.Sprintf("account %s was not found", address), http.StatusBadRequest)
		return
	}
	signature, err := c.keyManagerClient.SignTypedData(ctx, c.storeName, address, &qkmtypes.SignTypedDataRequest{
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
// @Accept json
// @Param request body qkmtypes.VerifyTypedDataRequest{domainSeparator=qkmtypes.DomainSeparator} true "Typed data to sign"
// @Success 204
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 422 {object} httputil.ErrorResponse "Invalid parameters"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/verify-typed-data [post]
func (c *AccountsController) verifyTypedDataSignature(rw http.ResponseWriter, request *http.Request) {
	verifyRequest := &qkmtypes.VerifyTypedDataRequest{}
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
// @Accept json
// @Param request body qkmtypes.VerifyRequest true "signature and message to verify"
// @Success 204
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 422 {object} httputil.ErrorResponse "Invalid parameters"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/verify-message [post]
func (c *AccountsController) verifyMessageSignature(rw http.ResponseWriter, request *http.Request) {
	verifyRequest := &qkmtypes.VerifyRequest{}
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
