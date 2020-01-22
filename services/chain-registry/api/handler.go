package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	"github.com/containous/traefik/v2/pkg/config/runtime"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ChainRegistryStore
}

func NewHandler(store types.ChainRegistryStore) *Handler {
	return &Handler{
		store: store,
	}
}

// Add internal routes to router
func (h *Handler) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/nodes").HandlerFunc(h.getNodes)
	router.Methods(http.MethodGet).Path("/nodes/{nodeID}").HandlerFunc(h.getNodeByID)
	router.Methods(http.MethodGet).Path("/{tenantID}/nodes").HandlerFunc(h.getNodesByTenantID)
	router.Methods(http.MethodGet).Path("/{tenantID}/nodes/{nodeID}").HandlerFunc(h.getNodeByTenantIDAndNodeID)

	router.Methods(http.MethodPost).Path("/nodes").HandlerFunc(h.postNode)
	router.Methods(http.MethodPost).Path("/{tenantID}/nodes").HandlerFunc(h.postNode)

	router.Methods(http.MethodPatch).Path("/{tenantID}/nodes/{nodeName}").HandlerFunc(h.patchNodeByName)
	router.Methods(http.MethodPatch).Path("/{tenantID}/nodes/{nodeName}/block-position").HandlerFunc(h.patchBlockPositionByName)
	router.Methods(http.MethodPatch).Path("/nodes/{nodeID}").HandlerFunc(h.patchNodeByID)
	router.Methods(http.MethodPatch).Path("/nodes/{nodeID}/block-position").HandlerFunc(h.patchBlockPositionByID)

	router.Methods(http.MethodDelete).Path("/{tenantID}/nodes/{nodeName}").HandlerFunc(h.deleteNodeByName)
	router.Methods(http.MethodDelete).Path("/nodes/{nodeID}").HandlerFunc(h.deleteNodeByID)
}

type Builder func(config *runtime.Configuration) http.Handler

type apiError struct {
	Message string `json:"message"`
}

func handleChainRegistryStoreError(rw http.ResponseWriter, err error) {
	if errors.IsNotFoundError(err) {
		writeError(rw, err.Error(), http.StatusNotFound)
	} else if err != nil {
		writeError(rw, err.Error(), http.StatusInternalServerError)
	}
}

func writeError(rw http.ResponseWriter, msg string, code int) {
	data, _ := json.Marshal(apiError{Message: msg})
	http.Error(rw, string(data), code)
}

type NodeRequest struct {
	Name                    string   `json:"name,omitempty"`
	URLs                    []string `json:"urls,omitempty" sql:"urls,array"`
	ListenerDepth           uint64   `json:"listenerDepth,omitempty"`
	ListenerBlockPosition   int64    `json:"listenerBlockPosition,string,omitempty"`
	ListenerFromBlock       int64    `json:"listenerFromBlock,string,omitempty"`
	ListenerBackOffDuration string   `json:"listenerBackOffDuration,omitempty"`
}

func UnmarshalBody(body io.ReadCloser, req interface{}) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields() // Force errors if unknown fields
	err := dec.Decode(req)
	if err != nil {
		return err
	}
	return nil
}

func UnmarshalNodeRequestBody(body io.ReadCloser) (*NodeRequest, error) {
	nodeRequest := &NodeRequest{}

	err := UnmarshalBody(body, nodeRequest)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(component)
	}

	// Check uniqueness of each urls
	keys := make(map[string]bool)
	for _, url := range nodeRequest.URLs {
		_, err := neturl.ParseRequestURI(url)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(component)
		}

		if _, value := keys[url]; value {
			return nil, errors.FromError(fmt.Errorf("cannot have twice the same url - got at least two times %s", url)).ExtendComponent(component)
		}
		keys[url] = true
	}

	return nodeRequest, nil
}
