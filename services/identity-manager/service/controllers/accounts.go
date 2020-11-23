package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/identitymanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/identity-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/identity-manager/service/formatters"
)

type IdentitiesController struct {
	accountUCs usecases.AccountUseCases
}

func NewIdentitiesController(accountUCs usecases.AccountUseCases) *IdentitiesController {
	return &IdentitiesController{
		accountUCs,
	}
}

// Add routes to router
func (c *IdentitiesController) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/accounts").HandlerFunc(c.search)
	router.Methods(http.MethodPost).Path("/accounts").HandlerFunc(c.create)
	router.Methods(http.MethodPost).Path("/accounts/import").HandlerFunc(c.importKey)
	router.Methods(http.MethodGet).Path("/accounts/{address}").HandlerFunc(c.getOne)
	router.Methods(http.MethodPatch).Path("/accounts/{address}").HandlerFunc(c.update)
	router.Methods(http.MethodPost).Path("/accounts/{address}/sign").HandlerFunc(c.signPayload)
}

// @Summary Creates a new Account
// @Description Creates a new Account
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body identitymanager.CreateAccountRequest true "Account creation request"
// @Success 200 {object} identitymanager.AccountResponse "Account object"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts [post]
func (c *IdentitiesController) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	req := &identitymanager.CreateAccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, req)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	acc := formatters.FormatCreateAccountRequest(req)
	acc, err = c.accountUCs.CreateAccount().Execute(ctx, acc, "", req.Chain, multitenancy.TenantIDFromContext(ctx))
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
// @Success 200 {object} identitymanager.AccountResponse "Account found"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/{address} [get]
func (c *IdentitiesController) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	address, err := utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	acc, err := c.accountUCs.GetAccount().Execute(ctx, address, multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(acc))
}

// @Summary Search accounts by provided filters
// @Description Get a list of filtered identities
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param aliases query []string false "List of account aliases" collectionFormat(csv)
// @Success 200 {array} identitymanager.AccountResponse "List of identities found"
// @Failure 400 {object} httputil.ErrorResponse "Invalid filter in the request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts [get]
func (c *IdentitiesController) search(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	filters, err := formatters.FormatAccountFilterRequest(request)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	accs, err := c.accountUCs.SearchAccounts().Execute(ctx, filters, multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	var response []*identitymanager.AccountResponse
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
// @Param request body identitymanager.ImportAccountRequest true "Account creation request"
// @Success 200 {object} identitymanager.AccountResponse "Account object"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 422 {object} httputil.ErrorResponse "Unprocessable entity"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 405 {object} httputil.ErrorResponse "Not allowed"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/import [post]
func (c *IdentitiesController) importKey(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	req := &identitymanager.ImportAccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, req)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	accResp := formatters.FormatImportAccountRequest(req)
	accResp, err = c.accountUCs.CreateAccount().Execute(ctx, accResp, req.PrivateKey, req.Chain, multitenancy.TenantIDFromContext(ctx))
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
// @Param request body identitymanager.UpdateAccountRequest{} true "Account update request"
// @Success 200 {object} identitymanager.AccountResponse "Account found"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/{address} [patch]
func (c *IdentitiesController) update(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	accRequest := &identitymanager.UpdateAccountRequest{}
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

	accRes, err := c.accountUCs.UpdateAccount().Execute(ctx, acc, multitenancy.AllowedTenantsFromContext(ctx))

	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatAccountResponse(accRes))
}

// @Summary Sign arbitrary data
// @Description Sign sent data using provided address
// @Accept json
// @Produce text/plain
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param address path string true "selected account address"
// @Success 200 {object} identitymanager.SignPayloadRequest "Data signature"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /accounts/{address}/sign [post]
func (c *IdentitiesController) signPayload(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	payloadRequest := &identitymanager.SignPayloadRequest{}
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

	signature, err := c.accountUCs.SignPayload().Execute(ctx, address, payloadRequest.Data, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(signature))
}
