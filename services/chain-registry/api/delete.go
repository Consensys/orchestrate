package api

import (
	"encoding/json"
	"net/http"

	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	"github.com/gorilla/mux"
)

type deleteResponse struct{}

func (h Handler) deleteChainByUUID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	err := h.store.DeleteChainByUUID(request.Context(), mux.Vars(request)["uuid"])
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&deleteResponse{})
}

func (h Handler) deleteChainByName(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chain := &models.Chain{
		Name:     mux.Vars(request)["name"],
		TenantID: mux.Vars(request)["tenantID"],
	}

	err := h.store.DeleteChainByName(request.Context(), chain)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&deleteResponse{})
}
