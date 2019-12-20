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

var (
	nodeIDPath             = "nodeID"
	tenantIDPath           = "tenantID"
	nodeNamePath           = "nodeName"
	getNodesPath           = "/api/nodes"
	getNodesByTenantIDPath = fmt.Sprintf("/api/nodes/{%s}", tenantIDPath)
	getNodeByNamePath      = fmt.Sprintf("/api/nodes/{%s}/{%s}", tenantIDPath, nodeNamePath)
	getNodeByIDPath        = fmt.Sprintf("/api/node/{%s}", nodeIDPath)
	postNodePath           = fmt.Sprintf("/api/nodes/{%s}", tenantIDPath)
	patchNodeByNamePath    = fmt.Sprintf("/api/nodes/{%s}/{%s}", tenantIDPath, nodeNamePath)
	patchNodeByIDPath      = fmt.Sprintf("/api/node/{%s}", nodeIDPath)
	deleteNodeByNamePath   = fmt.Sprintf("/api/nodes/{%s}/{%s}", tenantIDPath, nodeNamePath)
	deleteNodeByIDPath     = fmt.Sprintf("/api/node/{%s}", nodeIDPath)
)

// Add internal routes to router
func (h *Handler) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path(getNodesPath).HandlerFunc(h.getNodes)
	router.Methods(http.MethodGet).Path(getNodesByTenantIDPath).HandlerFunc(h.getNodesByTenantID)
	router.Methods(http.MethodGet).Path(getNodeByNamePath).HandlerFunc(h.getNodeByName)
	router.Methods(http.MethodGet).Path(getNodeByIDPath).HandlerFunc(h.getNodeByID)

	router.Methods(http.MethodPost).Path(postNodePath).HandlerFunc(h.postNode)

	router.Methods(http.MethodPatch).Path(patchNodeByNamePath).HandlerFunc(h.patchNodeByName)
	router.Methods(http.MethodPatch).Path(patchNodeByIDPath).HandlerFunc(h.patchNodeByID)

	router.Methods(http.MethodDelete).Path(deleteNodeByNamePath).HandlerFunc(h.deleteNodeByName)
	router.Methods(http.MethodDelete).Path(deleteNodeByIDPath).HandlerFunc(h.deleteNodeByID)
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
	data, err := json.Marshal(apiError{Message: msg})
	if err != nil {
		http.Error(rw, msg, code)
		return
	}

	http.Error(rw, string(data), code)
}

func UnmarshalNodeRequestBody(body io.ReadCloser) (*NodeRequest, error) {
	nodeRequest := &NodeRequest{}

	// Unmarshal body
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields() // Force errors if unknown fields
	err := dec.Decode(&nodeRequest)
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
