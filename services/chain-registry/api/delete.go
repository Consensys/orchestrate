package api

import (
	"encoding/json"
	"net/http"

	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	"github.com/gorilla/mux"
)

type deleteResponse struct{}

func (h Handler) deleteNodeByID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	nodeID := mux.Vars(request)[nodeIDPath]

	err := h.store.DeleteNodeByID(request.Context(), nodeID)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&deleteResponse{})
}

func (h Handler) deleteNodeByName(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	node := &models.Node{
		Name: mux.Vars(request)[nodeNamePath],
		// TODO: replace tenantID when extract token
		TenantID: mux.Vars(request)[tenantIDPath],
	}

	err := h.store.DeleteNodeByName(request.Context(), node)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&deleteResponse{})
}
