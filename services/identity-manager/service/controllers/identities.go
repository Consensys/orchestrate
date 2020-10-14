package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/multitenancy"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/identitymanager"
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
}

// @Summary Creates a new Identity
// @Description Creates a new Identity
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body identity.CreateIdentityRequest{} "Identity creation request"
// @Success 200 {object} identity.IdentityResponse{} "Identity object"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /identities [post]
func (c *IdentitiesController) create(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	req := &types.CreateIdentityRequest{}
	err := jsonutils.UnmarshalBody(request.Body, req)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	iden := formatters.FormatCreateIdentityRequest(req)
	iden, err = c.identityUCs.CreateIdentity().Execute(ctx, iden, multitenancy.TenantIDFromContext(ctx))
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatIdentityResponse(iden))
}
