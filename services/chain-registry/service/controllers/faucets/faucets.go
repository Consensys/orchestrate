package faucets

import (
	"net/http"

	"github.com/gorilla/mux"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/use-cases"
)

type Controller interface {
	Append(router *mux.Router)
}

type controller struct {
	getFaucetsUC   usecases.GetFaucets
	getFaucetUC    usecases.GetFaucet
	registerFaucet usecases.RegisterFaucet
	deleteFaucet   usecases.DeleteFaucet
	updateFaucet   usecases.UpdateFaucet
}

func NewController(
	getFaucetsUC usecases.GetFaucets,
	getFaucetUC usecases.GetFaucet,
	registerFaucet usecases.RegisterFaucet,
	deleteFaucet usecases.DeleteFaucet,
	updateFaucet usecases.UpdateFaucet,
) Controller {
	return &controller{
		getFaucetsUC:   getFaucetsUC,
		getFaucetUC:    getFaucetUC,
		registerFaucet: registerFaucet,
		deleteFaucet:   deleteFaucet,
		updateFaucet:   updateFaucet,
	}
}

// Add routes to router
func (h *controller) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/faucets").HandlerFunc(h.GetFaucets)
	router.Methods(http.MethodGet).Path("/faucets/{uuid}").HandlerFunc(h.GetFaucet)
	router.Methods(http.MethodPost).Path("/faucets").HandlerFunc(h.PostFaucet)
	router.Methods(http.MethodPatch).Path("/faucets/{uuid}").HandlerFunc(h.PatchFaucet)
	router.Methods(http.MethodDelete).Path("/faucets/{uuid}").HandlerFunc(h.DeleteFaucet)
}
