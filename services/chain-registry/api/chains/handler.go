package chains

import (
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
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

	router.Methods(http.MethodPost).Path("/chains").HandlerFunc(h.postChain)

	router.Methods(http.MethodPatch).Path("/chains/{uuid}").HandlerFunc(h.patchChainByUUID)

	router.Methods(http.MethodDelete).Path("/chains/{uuid}").HandlerFunc(h.deleteChainByUUID)
}

type Listener struct {
	Depth             *uint64 `json:"depth,omitempty"`
	BlockPosition     *int64  `json:"blockPosition,string,omitempty"`
	BackOffDuration   *string `json:"backOffDuration,omitempty"`
	ExternalTxEnabled *bool   `json:"externalTxEnabled,omitempty"`
}
