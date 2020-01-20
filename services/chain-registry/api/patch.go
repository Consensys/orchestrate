package api

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	"github.com/gorilla/mux"
	models "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"
)

type PatchBlockPositionRequest struct {
	BlockPosition int64 `json:"blockPosition,string,omitempty"`
}

type patchResponse struct{}

func (h Handler) patchNodeByName(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	nodeRequest, err := UnmarshalNodeRequestBody(request.Body)
	if err != nil {
		writeError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	node := &models.Node{
		Name: mux.Vars(request)["nodeName"],
		// TODO: replace tenantID when extract token
		TenantID:                mux.Vars(request)["tenantID"],
		URLs:                    nodeRequest.URLs,
		ListenerDepth:           nodeRequest.ListenerDepth,
		ListenerBlockPosition:   nodeRequest.ListenerBlockPosition,
		ListenerFromBlock:       nodeRequest.ListenerFromBlock,
		ListenerBackOffDuration: nodeRequest.ListenerBackOffDuration,
	}

	err = h.store.UpdateNodeByName(request.Context(), node)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&patchResponse{})
}

func (h Handler) patchBlockPositionByName(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	p := &PatchBlockPositionRequest{}
	err := UnmarshalBody(request.Body, p)
	if err != nil {
		writeError(rw, errors.FromError(err).ExtendComponent(component).Error(), http.StatusBadRequest)
		return
	}

	err = h.store.UpdateBlockPositionByName(request.Context(), mux.Vars(request)["nodeName"], mux.Vars(request)["tenantID"], p.BlockPosition)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&patchResponse{})
}

func (h Handler) patchNodeByID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	nodeID := mux.Vars(request)["nodeID"]

	nodeRequest, err := UnmarshalNodeRequestBody(request.Body)
	if err != nil {
		writeError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	node := &models.Node{
		ID:                      nodeID,
		Name:                    nodeRequest.Name,
		URLs:                    nodeRequest.URLs,
		ListenerDepth:           nodeRequest.ListenerDepth,
		ListenerBlockPosition:   nodeRequest.ListenerBlockPosition,
		ListenerFromBlock:       nodeRequest.ListenerFromBlock,
		ListenerBackOffDuration: nodeRequest.ListenerBackOffDuration,
	}

	err = h.store.UpdateNodeByID(request.Context(), node)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&patchResponse{})
}

func (h Handler) patchBlockPositionByID(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	nodeID := mux.Vars(request)["nodeID"]

	p := &PatchBlockPositionRequest{}
	err := UnmarshalBody(request.Body, p)
	if err != nil {
		writeError(rw, errors.FromError(err).ExtendComponent(component).Error(), http.StatusBadRequest)
		return
	}

	err = h.store.UpdateBlockPositionByID(request.Context(), nodeID, p.BlockPosition)
	if err != nil {
		handleChainRegistryStoreError(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(&patchResponse{})
}
