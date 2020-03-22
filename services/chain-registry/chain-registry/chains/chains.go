package chains

import (
	"context"
	"net/http"

	"github.com/containous/traefik/v2/pkg/log"
	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient"
	ethclientutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/ethclient/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store"
)

type Chains struct {
	store store.ChainStore
	ec    ethclient.ChainLedgerReader
}

func New(s store.ChainStore, ec ethclient.ChainLedgerReader) *Chains {
	return &Chains{
		store: s,
		ec:    ec,
	}
}

// Add routes to router
func (h *Chains) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/chains").HandlerFunc(h.GetChains)
	router.Methods(http.MethodGet).Path("/chains/{uuid}").HandlerFunc(h.GetChain)
	router.Methods(http.MethodPost).Path("/chains").HandlerFunc(h.PostChain)
	router.Methods(http.MethodPatch).Path("/chains/{uuid}").HandlerFunc(h.PatchChain)
	router.Methods(http.MethodDelete).Path("/chains/{uuid}").HandlerFunc(h.DeleteChain)
}

func (h *Chains) getChainTip(ctx context.Context, urls []string) (uint64, error) {
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
