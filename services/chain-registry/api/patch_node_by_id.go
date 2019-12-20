package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type patchNodeByIDResponse struct{}

func (h Handler) patchNodeByID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	nodeID := mux.Vars(request)[nodeIDPath]

	nodeRequest, err := UnmarshalNodeRequestBody(request.Body)
	if err != nil {
		writeError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	node := &models.Node{
		ID:                      nodeID,
		Name:                    nodeRequest.Name,
		URLs:                    nodeRequest.URLs,
		ListenerDepth:           nodeRequest.ListenerDepth,
		ListenerBlockPosition:   nodeRequest.ListenerBlockPosition,
		ListenerFromBlock:       nodeRequest.ListenerFromBlock,
		ListenerBackOffDuration: nodeRequest.ListenerBackOffDuration,
	}

	err = h.store.UpdateNodeByID(request.Context(), node)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&patchNodeByIDResponse{})
}
