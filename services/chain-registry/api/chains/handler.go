package chains

import (
	"context"
	"net/http"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient"
	ethclientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type Chain = types.Chain

type Handler struct {
	store types.ChainStore
	ec    ethclient.ChainLedgerReader
}

func NewHandler(store types.ChainStore, ec ethclient.ChainLedgerReader) *Handler {
	return &Handler{
		store: store,
		ec:    ec,
	}
}

// Add routes to router
func (h *Handler) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/chains").HandlerFunc(h.getChains)
	router.Methods(http.MethodGet).Path("/chains/{uuid}").HandlerFunc(h.getChainByUUID)

	router.Methods(http.MethodPost).Path("/chains").HandlerFunc(h.postChain)

	router.Methods(http.MethodPatch).Path("/chains/{uuid}").HandlerFunc(h.patchChainByUUID)

	router.Methods(http.MethodDelete).Path("/chains/{uuid}").HandlerFunc(h.deleteChainByUUID)
}

func (h *Handler) getChainTip(ctx context.Context, urls []string) (uint64, error) {
	var tip uint64

	// All URLs must be valid and we return the head of the latest one
	for _, url := range urls {
		head, err := h.ec.HeaderByNumber(ethclientutils.RetryNotFoundError(ctx, true), url, nil)
		if err != nil {
			log.FromContext(ctx).WithError(err).Errorf("failed to fetch chain tip for URL %s", url)
			return 0, err
		}

		tip = head.Number.Uint64()
	}

	return tip, nil
}
