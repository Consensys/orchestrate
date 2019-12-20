package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type deleteNodeByIDResponse struct{}

func (h Handler) deleteNodeByID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	nodeID := mux.Vars(request)[nodeIDPath]

	err := h.store.DeleteNodeByID(request.Context(), nodeID)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&deleteNodeByIDResponse{})
}
