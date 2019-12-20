package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (h Handler) getNodeByID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	nodeID, err := strconv.Atoi(mux.Vars(request)[nodeIDPath])
	if err != nil {
		writeError(rw, "invalid ID format", http.StatusBadRequest)
		return
	}

	node, err := h.store.GetNodeByID(request.Context(), nodeID)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(node)
}
