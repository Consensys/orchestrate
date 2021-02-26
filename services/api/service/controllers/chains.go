package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"

	jsonutils "github.com/ConsenSys/orchestrate/pkg/encoding/json"
	"github.com/ConsenSys/orchestrate/pkg/multitenancy"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/ConsenSys/orchestrate/pkg/types/api"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/ConsenSys/orchestrate/services/api/service/formatters"
	"github.com/gorilla/mux"
)

// Hack for swagger generation
var _ entities.PrivateTxManager

type ChainsController struct {
	ucs usecases.ChainUseCases
}

func NewChainsController(chainUCs usecases.ChainUseCases) *ChainsController {
	return &ChainsController{ucs: chainUCs}
}

func (c *ChainsController) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/chains").HandlerFunc(c.search)
	router.Methods(http.MethodGet).Path("/chains/{uuid}").HandlerFunc(c.getOne)
	router.Methods(http.MethodPost).Path("/chains").HandlerFunc(c.register)
	router.Methods(http.MethodPatch).Path("/chains/{uuid}").HandlerFunc(c.update)
	router.Methods(http.MethodDelete).Path("/chains/{uuid}").HandlerFunc(c.delete)
}

// @Summary Retrieves a list of all registered chains
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200 {array} api.ChainResponse{privateTxManager=entities.PrivateTxManager}
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /chains [get]
func (c *ChainsController) search(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	filters, err := formatters.FormatChainFiltersRequest(request)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	chains, err := c.ucs.SearchChains().Execute(ctx, filters, multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	response := []*api.ChainResponse{}
	for _, chain := range chains {
		response = append(response, formatters.FormatChainResponse(chain))
	}

	_ = json.NewEncoder(rw).Encode(response)
}

// @Summary Retrieves a chain by ID
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the chain"
// @Success 200 {object} api.ChainResponse{privateTxManager=entities.PrivateTxManager}
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Chain not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /chains/{uuid} [get]
func (c *ChainsController) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	chain, err := c.ucs.GetChain().Execute(ctx, mux.Vars(request)["uuid"], multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatChainResponse(chain))
}

// @Summary Updates a chain by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the chain"
// @Param request body api.UpdateChainRequest{listener=api.UpdateListenerRequest,privateTxManager=api.PrivateTxManagerRequest} true "Chain update request"
// @Success 200 {object} api.ChainResponse{privateTxManager=entities.PrivateTxManager}
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Chain not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /chains/{uuid} [patch]
func (c *ChainsController) update(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	chainRequest := &api.UpdateChainRequest{}
	err := jsonutils.UnmarshalBody(request.Body, chainRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	uuid := mux.Vars(request)["uuid"]
	allowedTenants := multitenancy.AllowedTenantsFromContext(ctx)
	chain, err := c.ucs.UpdateChain().Execute(ctx, formatters.FormatUpdateChainRequest(chainRequest, uuid), allowedTenants)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatChainResponse(chain))
}

// @Summary Registers a new chain
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body api.RegisterChainRequest{listener=api.RegisterListenerRequest,privateTxManager=api.PrivateTxManagerRequest} true "Chain registration request"
// @Success 200 {object} api.ChainResponse{privateTxManager=entities.PrivateTxManager}
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /chains [post]
func (c *ChainsController) register(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	chainRequest := &api.RegisterChainRequest{}
	err := jsonutils.UnmarshalBody(request.Body, chainRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	fromLatest := chainRequest.Listener.FromBlock == "" || chainRequest.Listener.FromBlock == "latest"
	tenantID := multitenancy.TenantIDFromContext(ctx)
	chain, err := formatters.FormatRegisterChainRequest(chainRequest, tenantID, fromLatest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	chain, err = c.ucs.RegisterChain().Execute(ctx, chain, fromLatest)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatChainResponse(chain))
}

// @Summary Deletes a chain by ID
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the chain"
// @Success 204
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Chain not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /chains/{uuid} [delete]
func (c *ChainsController) delete(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	uuid := mux.Vars(request)["uuid"]
	tenants := multitenancy.AllowedTenantsFromContext(ctx)

	err := c.ucs.DeleteChain().Execute(ctx, uuid, tenants)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
