package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type postNodeResponse struct {
	ID string `json:"id,omitempty"`
}

func (h Handler) postNode(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	nodeRegisterRequest, err := UnmarshalNodeRegisterRequestBody(request.Body)
	if err != nil {
		writeError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	node := &models.Node{
		Name: nodeRegisterRequest.Name,
		//TODO: replace tenantID by the one extracted in the token when ready
		TenantID:                mux.Vars(request)["tenantID"],
		URLs:                    nodeRegisterRequest.URLs,
		ListenerDepth:           nodeRegisterRequest.ListenerDepth,
		ListenerBlockPosition:   nodeRegisterRequest.ListenerFromBlock,
		ListenerFromBlock:       nodeRegisterRequest.ListenerFromBlock,
		ListenerBackOffDuration: nodeRegisterRequest.ListenerBackOffDuration,
	}

	err = h.store.RegisterNode(request.Context(), node)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&postNodeResponse{
		ID: node.ID,
	})
}
