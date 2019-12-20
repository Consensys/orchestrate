package api

import (
	"encoding/json"
	"net/http"
)

func (h Handler) getNodes(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	node, err := h.store.GetNodes(request.Context())
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(node)
}
