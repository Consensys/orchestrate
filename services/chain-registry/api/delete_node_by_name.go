package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type deleteNodeByNameResponse struct{}

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

	_ = json.NewEncoder(rw).Encode(&deleteNodeByNameResponse{})
}
