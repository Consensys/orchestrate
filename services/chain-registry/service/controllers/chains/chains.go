package chains

import (
	"net/http"

	"github.com/gorilla/mux"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/use-cases/chains"
)

type Controller interface {
	Append(router *mux.Router)
}

type controller struct {
	getChainsUC     usecases.GetChains
	getChainUC      usecases.GetChain
	registerChainUC usecases.RegisterChain
	deleteChainUC   usecases.DeleteChain
	updateChainUC   usecases.UpdateChain
}

func NewController(
	getChainsUC usecases.GetChains,
	getChainUC usecases.GetChain,
	registerChainUC usecases.RegisterChain,
	deleteChainUC usecases.DeleteChain,
	updateChainUC usecases.UpdateChain,
) Controller {
	return &controller{
		getChainsUC:     getChainsUC,
		getChainUC:      getChainUC,
		registerChainUC: registerChainUC,
		deleteChainUC:   deleteChainUC,
		updateChainUC:   updateChainUC,
	}
}

// Add routes to router
func (h *controller) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/chains").HandlerFunc(h.GetChains)
	router.Methods(http.MethodGet).Path("/chains/{uuid}").HandlerFunc(h.GetChain)
	router.Methods(http.MethodPost).Path("/chains").HandlerFunc(h.PostChain)
	router.Methods(http.MethodPatch).Path("/chains/{uuid}").HandlerFunc(h.PatchChain)
	router.Methods(http.MethodDelete).Path("/chains/{uuid}").HandlerFunc(h.DeleteChain)
}
