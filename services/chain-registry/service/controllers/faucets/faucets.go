package faucets

import (
	"net/http"

	"github.com/gorilla/mux"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/chain-registry/use-cases/faucets"
)

type Controller interface {
	Append(router *mux.Router)
}

type controller struct {
	getFaucetsUC      usecases.GetFaucets
	getFaucetUC       usecases.GetFaucet
	registerFaucetUC  usecases.RegisterFaucet
	deleteFaucetUC    usecases.DeleteFaucet
	updateFaucetUC    usecases.UpdateFaucet
	faucetCandidateUC usecases.FaucetCandidate
}

func NewController(
	getFaucetsUC usecases.GetFaucets,
	getFaucetUC usecases.GetFaucet,
	registerFaucetUC usecases.RegisterFaucet,
	deleteFaucetUC usecases.DeleteFaucet,
	updateFaucetUC usecases.UpdateFaucet,
	faucetCandidateUC usecases.FaucetCandidate,
) Controller {
	return &controller{
		getFaucetsUC:      getFaucetsUC,
		getFaucetUC:       getFaucetUC,
		registerFaucetUC:  registerFaucetUC,
		deleteFaucetUC:    deleteFaucetUC,
		updateFaucetUC:    updateFaucetUC,
		faucetCandidateUC: faucetCandidateUC,
	}
}

// Add routes to router
func (h *controller) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/faucets").HandlerFunc(h.GetFaucets)
	router.Methods(http.MethodGet).Path("/faucets/candidate").HandlerFunc(h.GetFaucetCandidate)
	router.Methods(http.MethodGet).Path("/faucets/{uuid}").HandlerFunc(h.GetFaucet)
	router.Methods(http.MethodPost).Path("/faucets").HandlerFunc(h.PostFaucet)
	router.Methods(http.MethodPatch).Path("/faucets/{uuid}").HandlerFunc(h.PatchFaucet)
	router.Methods(http.MethodDelete).Path("/faucets/{uuid}").HandlerFunc(h.DeleteFaucet)
}
