package api

import (
	"encoding/json"
	"html"
	"net/http"

	"github.com/gorilla/mux"
)

func (h Handler) getChains(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	filters := make(map[string]string)
	for k := range request.URL.Query() {
		key := html.EscapeString(k)
		filters[key] = html.EscapeString(request.URL.Query().Get(k))
	}

	chains, err := h.store.GetChains(request.Context(), filters)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chains)
}

func (h Handler) getChainByUUID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chain, err := h.store.GetChainByUUID(request.Context(), mux.Vars(request)["uuid"])
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chain)
}

func (h Handler) getChainsByTenantID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	filters := make(map[string]string)
	for k := range request.URL.Query() {
		key := html.EscapeString(k)
		filters[key] = html.EscapeString(request.URL.Query().Get(k))
	}

	chains, err := h.store.GetChainsByTenantID(request.Context(), mux.Vars(request)["tenantID"], filters)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chains)
}

func (h Handler) getChainByTenantIDAndName(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	chain, err := h.store.GetChainByTenantIDAndName(request.Context(), mux.Vars(request)["tenantID"], mux.Vars(request)["name"])
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(chain)
}
