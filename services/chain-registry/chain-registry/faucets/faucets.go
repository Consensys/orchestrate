package faucets

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

type Faucets struct {
	store store.FaucetStore
}

func New(s store.FaucetStore) *Faucets {
	return &Faucets{
		store: s,
	}
}

// Add routes to router
func (h *Faucets) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/faucets").HandlerFunc(h.GetFaucets)
	router.Methods(http.MethodGet).Path("/faucets/{uuid}").HandlerFunc(h.GetFaucet)

	router.Methods(http.MethodPost).Path("/faucets").HandlerFunc(h.PostFaucet)

	router.Methods(http.MethodPatch).Path("/faucets/{uuid}").HandlerFunc(h.PatchFaucet)
	router.Methods(http.MethodDelete).Path("/faucets/{uuid}").HandlerFunc(h.DeleteFaucet)
}
