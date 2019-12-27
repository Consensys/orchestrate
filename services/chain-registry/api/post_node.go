package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type NodeRequest struct {
	Name                    string   `json:"name,omitempty"`
	URLs                    []string `json:"urls,omitempty" sql:"urls,array"`
	ListenerDepth           uint64   `json:"listenerDepth,omitempty"`
	ListenerBlockPosition   int64    `json:"listenerBlockPosition,string,omitempty"`
	ListenerFromBlock       int64    `json:"listenerFromBlock,string,omitempty"`
	ListenerBackOffDuration string   `json:"listenerBackOffDuration,omitempty"`
}

type postNodeResponse struct {
	ID string `json:"id,omitempty"`
}

func (h Handler) postNode(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	nodeRequest, err := UnmarshalNodeRequestBody(request.Body)
	if err != nil {
		writeError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	node := &models.Node{
		Name: nodeRequest.Name,
		//TODO: replace tenantID by the one extracted in the token when ready
		TenantID:                mux.Vars(request)["tenantID"],
		URLs:                    nodeRequest.URLs,
		ListenerDepth:           nodeRequest.ListenerDepth,
		ListenerBlockPosition:   nodeRequest.ListenerBlockPosition,
		ListenerFromBlock:       nodeRequest.ListenerFromBlock,
		ListenerBackOffDuration: nodeRequest.ListenerBackOffDuration,
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
