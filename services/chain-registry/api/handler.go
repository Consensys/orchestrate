package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	"github.com/containous/traefik/v2/pkg/config/runtime"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ChainRegistryStore
}

func NewHandler(store types.ChainRegistryStore) *Handler {
	return &Handler{
		store: store,
	}
}

// Add internal routes to router
func (h *Handler) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/chains").HandlerFunc(h.getChains)
	router.Methods(http.MethodGet).Path("/chains/{uuid}").HandlerFunc(h.getChainByUUID)
	router.Methods(http.MethodGet).Path("/{tenantID}/chains").HandlerFunc(h.getChainsByTenantID)
	router.Methods(http.MethodGet).Path("/{tenantID}/chains/{name}").HandlerFunc(h.getChainByTenantIDAndName)

	router.Methods(http.MethodPost).Path("/chains").HandlerFunc(h.postChain)
	router.Methods(http.MethodPost).Path("/{tenantID}/chains").HandlerFunc(h.postChain)

	router.Methods(http.MethodPatch).Path("/{tenantID}/chains/{name}").HandlerFunc(h.patchChainByName)
	router.Methods(http.MethodPatch).Path("/chains/{uuid}").HandlerFunc(h.patchChainByUUID)

	router.Methods(http.MethodDelete).Path("/{tenantID}/chains/{name}").HandlerFunc(h.deleteChainByName)
	router.Methods(http.MethodDelete).Path("/chains/{uuid}").HandlerFunc(h.deleteChainByUUID)
}

type Builder func(config *runtime.Configuration) http.Handler

type apiError struct {
	Message string `json:"message"`
}

func handleChainRegistryStoreError(rw http.ResponseWriter, err error) {
	switch {
	case errors.IsAlreadyExistsError(err):
		writeError(rw, err.Error(), http.StatusConflict)
	case errors.IsNotFoundError(err):
		writeError(rw, err.Error(), http.StatusNotFound)
	case err != nil:
		writeError(rw, err.Error(), http.StatusInternalServerError)
	}
}

func writeError(rw http.ResponseWriter, msg string, code int) {
	data, _ := json.Marshal(apiError{Message: msg})
	http.Error(rw, string(data), code)
}

type Listener struct {
	Depth           *uint64 `json:"depth,omitempty"`
	BlockPosition   *int64  `json:"blockPosition,string,omitempty"`
	BackOffDuration *string `json:"backOffDuration,omitempty"`
}

var validate = validator.New()

func UnmarshalBody(body io.Reader, req interface{}) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields() // Force errors if unknown fields
	err := dec.Decode(req)
	if err != nil {
		return errors.FromError(err).ExtendComponent(component)
	}

	err = validate.Struct(req)
	if err != nil {
		return errors.FromError(fmt.Errorf("invalid body")).ExtendComponent(component)
	}

	return nil
}
