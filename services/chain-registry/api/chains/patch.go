package chains

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/api/utils"
	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type PatchRequest struct {
	Name     string    `json:"name,omitempty"`
	URLs     []string  `json:"urls,omitempty" pg:"urls,array" validate:"unique,dive,url"`
	Listener *Listener `json:"listener,omitempty"`
}

type PatchResponse struct{}

func (h Handler) patchChainByName(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chainRequest := &PatchRequest{Listener: &Listener{}}
	err := utils.UnmarshalBody(request.Body, chainRequest)
	if err != nil {
		utils.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	chain := &models.Chain{
		Name:     mux.Vars(request)["name"],
		TenantID: mux.Vars(request)["tenantID"],
		URLs:     chainRequest.URLs,
	}
	if chainRequest.Listener != nil {
		chain.ListenerDepth = chainRequest.Listener.Depth
		chain.ListenerBlockPosition = chainRequest.Listener.BlockPosition
		chain.ListenerBackOffDuration = chainRequest.Listener.BackOffDuration
		chain.ListenerExternalTxEnabled = chainRequest.Listener.ExternalTxEnabled
	}

	err = h.store.UpdateChainByName(request.Context(), chain)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&PatchResponse{})
}

func (h Handler) patchChainByUUID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chainRequest := &PatchRequest{Listener: &Listener{}}
	err := utils.UnmarshalBody(request.Body, chainRequest)
	if err != nil {
		utils.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	chain := &models.Chain{
		UUID:                      mux.Vars(request)["uuid"],
		Name:                      chainRequest.Name,
		URLs:                      chainRequest.URLs,
		ListenerDepth:             chainRequest.Listener.Depth,
		ListenerBlockPosition:     chainRequest.Listener.BlockPosition,
		ListenerBackOffDuration:   chainRequest.Listener.BackOffDuration,
		ListenerExternalTxEnabled: chainRequest.Listener.ExternalTxEnabled,
	}

	err = h.store.UpdateChainByUUID(request.Context(), chain)
	if err != nil {
		utils.HandleStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&PatchResponse{})
}
