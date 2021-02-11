package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/service/formatters"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/client"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
)

type AccountsController struct {
	ucs              usecases.AccountUseCases
	keyManagerClient client.KeyManagerClient
}

func NewAccountsController(accountUCs usecases.AccountUseCases, keyManagerClient client.KeyManagerClient) *AccountsController {
	return &AccountsController{
		accountUCs,
		keyManagerClient,
	}
}

// Add routes to router
func (c *AccountsController) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/accounts").HandlerFunc(c.search)
	router.Methods(http.MethodPost).Path("/accounts").HandlerFunc(c.create)
	router.Methods(http.MethodPost).Path("/accounts/import").HandlerFunc(c.importKey)
	router.Methods(http.MethodGet).Path("/accounts/{address}").HandlerFunc(c.getOne)
	router.Methods(http.MethodPatch).Path("/accounts/{address}").HandlerFunc(c.update)
	router.Methods(http.MethodPost).Path("/accounts/{address}/sign").HandlerFunc(c.signPayload)
	router.Methods(http.MethodPost).Path("/accounts/{address}/sign-typed-data").HandlerFunc(c.signTypedData)
	router.Methods(http.MethodPost).Path("/accounts/verify-signature").HandlerFunc(c.verifySignature)
	router.Methods(http.MethodPost).Path("/accounts/verify-typed-data-signature").HandlerFunc(c.verifyTypedDataSignature)
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
	acc, err = c.ucs.CreateAccount().Execute(ctx, acc, "", req.Chain, multitenancy.TenantIDFromContext(ctx))
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

	accResp := formatters.FormatImportAccountRequest(req)
	accResp, err = c.ucs.CreateAccount().Execute(ctx, accResp, req.PrivateKey, req.Chain, multitenancy.TenantIDFromContext(ctx))
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

// @Summary Sign arbitrary data
// @Description Sign sent data using provided account
// @Accept json
// @Produce text/plain
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param address path string true "selected account address"
// @Success 200 {object} api.SignPayloadRequest "Data signature"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/{address}/sign [post]
func (c *AccountsController) signPayload(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	payloadRequest := &api.SignPayloadRequest{}
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

	signature, err := c.keyManagerClient.ETHSign(request.Context(), address, &keymanager.SignPayloadRequest{
		Namespace: multitenancy.TenantIDFromContext(ctx),
		Data:      payloadRequest.Data,
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
// @Param request body api.SignTypedDataRequest{domainSeparator=ethereum.DomainSeparator,types=map[string]ethereum.Type} true "Typed data to sign"
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

	signature, err := c.keyManagerClient.ETHSignTypedData(ctx, address, &ethereum.SignTypedDataRequest{
		Namespace:       multitenancy.TenantIDFromContext(ctx),
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
// @Param request body ethereum.VerifyTypedDataRequest{domainSeparator=ethereum.DomainSeparator} true "Typed data to sign"
// @Success 204
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 422 {object} httputil.ErrorResponse "Invalid parameters"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/verify-typed-data-signature [post]
func (c *AccountsController) verifyTypedDataSignature(rw http.ResponseWriter, request *http.Request) {
	verifyRequest := &ethereum.VerifyTypedDataRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	verifyRequest.Address, err = utils.ParseHexToMixedCaseEthAddress(verifyRequest.Address)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.keyManagerClient.ETHVerifyTypedDataSignature(request.Context(), verifyRequest)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Verifies the signature of a message
// @Description Verifies if a message has been signed by the Ethereum account passed as argument
// @Accept json
// @Param request body ethereum.VerifyPayloadRequest true "signature and message to verify"
// @Success 204
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 422 {object} httputil.ErrorResponse "Invalid parameters"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/verify-signature [post]
func (c *AccountsController) verifySignature(rw http.ResponseWriter, request *http.Request) {
	verifyRequest := &ethereum.VerifyPayloadRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	verifyRequest.Address, err = utils.ParseHexToMixedCaseEthAddress(verifyRequest.Address)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.keyManagerClient.ETHVerifySignature(request.Context(), verifyRequest)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
