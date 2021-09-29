package controllers

import (
	"encoding/json"
	"net/http"

	jsonutils "github.com/consensys/orchestrate/pkg/encoding/json"
	"github.com/consensys/orchestrate/pkg/multitenancy"
	"github.com/consensys/orchestrate/pkg/toolkit/app/http/httputil"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/services/api/service/formatters"

	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/gorilla/mux"
)

type FaucetsController struct {
	ucs usecases.FaucetUseCases
}

func NewFaucetsController(ucs usecases.FaucetUseCases) *FaucetsController {
	return &FaucetsController{ucs: ucs}
}

// Add routes to router
func (c *FaucetsController) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/faucets").HandlerFunc(c.search)
	router.Methods(http.MethodGet).Path("/faucets/{uuid}").HandlerFunc(c.getOne)
	router.Methods(http.MethodPost).Path("/faucets").HandlerFunc(c.register)
	router.Methods(http.MethodPatch).Path("/faucets/{uuid}").HandlerFunc(c.update)
	router.Methods(http.MethodDelete).Path("/faucets/{uuid}").HandlerFunc(c.delete)
}

// @Summary Retrieves a list of all registered faucets
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200 {array} api.FaucetResponse
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /faucets [get]
func (c *FaucetsController) search(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	filters, err := formatters.FormatFaucetFilters(request)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	faucets, err := c.ucs.SearchFaucets().Execute(ctx, filters, multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	response := []*api.FaucetResponse{}
	for _, faucet := range faucets {
		response = append(response, formatters.FormatFaucetResponse(faucet))
	}

	_ = json.NewEncoder(rw).Encode(response)
}

// @Summary Retrieves a faucet by ID
// @Produce json
// @Param uuid path string true "ID of the faucet"
// @Success 200 {object} api.FaucetResponse
// @Failure 404 {object} httputil.ErrorResponse "Faucet not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /faucets/{uuid} [get]
func (c *FaucetsController) getOne(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	faucet, err := c.ucs.GetFaucet().Execute(ctx, mux.Vars(request)["uuid"], multitenancy.AllowedTenantsFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatFaucetResponse(faucet))
}

// @Summary Registers a new faucet
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body api.RegisterFaucetRequest true "Faucet registration request"
// @Success 200 {object} api.FaucetResponse
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 422 {object} httputil.ErrorResponse "Unprocessable entity"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /faucets [post]
func (c *FaucetsController) register(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	faucetRequest := &api.RegisterFaucetRequest{}
	err := jsonutils.UnmarshalBody(request.Body, faucetRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	tenantID := multitenancy.TenantIDFromContext(ctx)
	faucet, err := c.ucs.RegisterFaucet().Execute(ctx, formatters.FormatRegisterFaucetRequest(faucetRequest, tenantID))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatFaucetResponse(faucet))
}

// @Summary Updates a faucet by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the faucet"
// @Param request body api.UpdateFaucetRequest true "Faucet update request"
// @Success 200 {object} api.FaucetResponse
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Faucet not found"
// @Failure 422 {object} httputil.ErrorResponse "Unprocessable entity"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /faucets/{uuid} [patch]
func (c *FaucetsController) update(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	faucetRequest := &api.UpdateFaucetRequest{}
	err := jsonutils.UnmarshalBody(request.Body, faucetRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	uuid := mux.Vars(request)["uuid"]
	tenantID := multitenancy.TenantIDFromContext(ctx)
	allowedTenants := multitenancy.AllowedTenantsFromContext(ctx)
	faucet, err := c.ucs.UpdateFaucet().Execute(ctx, formatters.FormatUpdateFaucetRequest(faucetRequest, uuid, tenantID), allowedTenants)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatFaucetResponse(faucet))
}

// @Summary Deletes a faucet by ID
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the faucet"
// @Success 204
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Faucet not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /faucets/{uuid} [delete]
func (c *FaucetsController) delete(rw http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	uuid := mux.Vars(request)["uuid"]
	tenants := multitenancy.AllowedTenantsFromContext(ctx)

	err := c.ucs.DeleteFaucet().Execute(ctx, uuid, tenants)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
