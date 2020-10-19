package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/identitymanager"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/identity-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/identity-manager/service/formatters"
)

type IdentitiesController struct {
	identityUCs usecases.IdentityUseCases
}

func NewIdentitiesController(identityUCs usecases.IdentityUseCases) *IdentitiesController {
	return &IdentitiesController{
		identityUCs,
	}
}

// Add routes to router
func (c *IdentitiesController) Append(router *mux.Router) {
	router.Methods(http.MethodPost).Path("/identities").HandlerFunc(c.create)
	router.Methods(http.MethodPost).Path("/identities/import").HandlerFunc(c.importKey)
}

// @Summary Creates a new Identity
// @Description Creates a new Identity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body identitymanager.CreateIdentityRequest true "Identity creation request"
// @Success 200 {object} identitymanager.IdentityResponse "Identity object"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /identities [post]
func (c *IdentitiesController) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	req := &identitymanager.CreateIdentityRequest{}
	err := jsonutils.UnmarshalBody(request.Body, req)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	iden := formatters.FormatCreateIdentityRequest(req)
	iden, err = c.identityUCs.CreateIdentity().Execute(ctx, iden, "", req.Chain, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatIdentityResponse(iden))
}

// @Summary Creates a new Identity by importing a private key
// @Description Creates a new Identity by importing a private key
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body identitymanager.ImportIdentityRequest true "Identity creation request"
// @Success 200 {object} identitymanager.IdentityResponse "Identity object"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 422 {object} httputil.ErrorResponse "Unprocessable entity"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 405 {object} httputil.ErrorResponse "Not allowed"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /identities/import [post]
func (c *IdentitiesController) importKey(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	req := &identitymanager.ImportIdentityRequest{}
	err := jsonutils.UnmarshalBody(request.Body, req)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	iden := formatters.FormatImportIdentityRequest(req)
	iden, err = c.identityUCs.CreateIdentity().Execute(ctx, iden, req.PrivateKey, req.Chain, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatIdentityResponse(iden))
}
