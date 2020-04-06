package chains

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/chain-registry/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/models"
)

// @Summary Updates a chain by ID
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param uuid path string true "ID of the chain"
// @Param request body PatchRequest true "Chain update request"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /chains/{uuid} [patch]
func (h *controller) PatchChain(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chain, err := parsePatchReqToChain(request)
	if err != nil {
		utils.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	chainUUID := mux.Vars(request)["uuid"]
	err = h.updateChainUC.Execute(request.Context(), chainUUID, "", chain)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	chain, err = h.getChainUC.Execute(request.Context(), chainUUID, "")
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chain)
}

func parsePatchReqToChain(request *http.Request) (*models.Chain, error) {
	chainRequest := &PatchRequest{
		Listener: &ListenerPatchRequest{},
	}

	err := jsonutils.UnmarshalBody(request.Body, chainRequest)
	if err != nil {
		return nil, err
	}

	chain := models.Chain{
		Name:                      chainRequest.Name,
		URLs:                      chainRequest.URLs,
		ListenerCurrentBlock:      chainRequest.Listener.CurrentBlock,
		ListenerBackOffDuration:   chainRequest.Listener.BackOffDuration,
		ListenerDepth:             chainRequest.Listener.Depth,
		ListenerExternalTxEnabled: chainRequest.Listener.ExternalTxEnabled,
	}

	if chainRequest.PrivateTxManager != nil {
		chain.PrivateTxManagers = []*models.PrivateTxManagerModel{
			{
				URL:  chainRequest.PrivateTxManager.URL,
				Type: chainRequest.PrivateTxManager.Type,
			},
		}
	}

	return &chain, nil
}
