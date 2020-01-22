package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (h Handler) getNodes(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	filters := make(map[string]string)
	for k := range request.URL.Query() {
		filters[k] = request.URL.Query().Get(k)
	}

	nodes, err := h.store.GetNodes(request.Context(), filters)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(nodes)
}

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

func (h Handler) getNodesByTenantID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	filters := make(map[string]string)
	for k := range request.URL.Query() {
		filters[k] = request.URL.Query().Get(k)
	}

	nodes, err := h.store.GetNodesByTenantID(request.Context(), mux.Vars(request)["tenantID"], filters)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(nodes)
}

func (h Handler) getNodeByTenantIDAndNodeID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	node, err := h.store.GetNodeByTenantIDAndNodeID(request.Context(), mux.Vars(request)["tenantID"], mux.Vars(request)["nodeID"])
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(node)
}
