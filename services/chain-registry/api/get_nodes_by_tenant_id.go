package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (h Handler) getNodesByTenantID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	nodes, err := h.store.GetNodesByTenantID(request.Context(), mux.Vars(request)[tenantIDPath])
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(nodes)
}
