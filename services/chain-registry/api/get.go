package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (h Handler) getNodeByID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	nodeID := mux.Vars(request)["nodeID"]

	node, err := h.store.GetNodeByID(request.Context(), nodeID)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(node)
}

func (h Handler) getNodeByName(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	node, err := h.store.GetNodeByName(request.Context(), mux.Vars(request)["tenantID"], mux.Vars(request)["nodeName"])
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(node)
}

func (h Handler) getNodes(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	nodes, err := h.store.GetNodes(request.Context())
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(nodes)
}

func (h Handler) getNodesByTenantID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	nodes, err := h.store.GetNodesByTenantID(request.Context(), mux.Vars(request)["tenantID"])
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(nodes)
}
