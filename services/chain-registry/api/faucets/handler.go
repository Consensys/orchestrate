package faucets

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type Faucet = types.Faucet

type Handler struct {
	store types.FaucetStore
}

func NewHandler(store types.FaucetStore) *Handler {
	return &Handler{
		store: store,
	}
}

// Add routes to router
func (h *Handler) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/faucets").HandlerFunc(h.getFaucets)
	router.Methods(http.MethodGet).Path("/faucets/{uuid}").HandlerFunc(h.getFaucetByUUID)

	router.Methods(http.MethodPost).Path("/faucets").HandlerFunc(h.postFaucet)

	router.Methods(http.MethodPatch).Path("/faucets/{uuid}").HandlerFunc(h.patchFaucetByUUID)
	router.Methods(http.MethodDelete).Path("/faucets/{uuid}").HandlerFunc(h.deleteFaucetByUUID)
}
