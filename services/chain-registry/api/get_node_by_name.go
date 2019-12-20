package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (h Handler) getNodeByName(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	node, err := h.store.GetNodeByName(request.Context(), mux.Vars(request)[tenantIDPath], mux.Vars(request)[nodeNamePath])
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(node)
}
