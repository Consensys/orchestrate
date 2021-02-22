package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/service/formatters"
)

var _ entities.Method
var _ entities.Event

type ContractsController struct {
	ucs usecases.ContractUseCases
}

func NewContractsController(contractUCs usecases.ContractUseCases) *ContractsController {
	return &ContractsController{
		ucs: contractUCs,
	}
}

// Add routes to router
func (c *ContractsController) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/contracts").HandlerFunc(c.getCatalog)
	router.Methods(http.MethodPost).Path("/contracts").HandlerFunc(c.register)
	router.Methods(http.MethodPost).Path("/contracts/accounts/{chain_id}/{address}").HandlerFunc(c.setCodeHash)
	router.Methods(http.MethodGet).Path("/contracts/accounts/{chain_id}/{address}/events").HandlerFunc(c.getEvents)
	router.Methods(http.MethodGet).Path("/contracts/{name}").HandlerFunc(c.getTags)
	router.Methods(http.MethodGet).Path("/contracts/{name}/{tag}").HandlerFunc(c.getContract)
	router.Methods(http.MethodGet).Path("/contracts/{name}/{tag}/method-signatures").HandlerFunc(c.getContractMethodSignatures)
}

// @Summary Returns a list of all registered contracts
// @Description Returns a list of all registered contracts
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200 {array} string "Registered contract List"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /contracts [get]
func (c *ContractsController) getCatalog(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	names, err := c.ucs.GetContractsCatalog().Execute(ctx)

	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	if names == nil {
		names = []string{}
	}
	_ = json.NewEncoder(rw).Encode(names)
}

// @Summary Register new solidity contract
// @Description Register new solidity contract in Orchestrate
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param request body api.RegisterContractRequest true "Contract register request"
// @Success 200 {object} api.ContractResponse{constructor=entities.Method,methods=[]entities.Method,events=[]entities.Event} "Contract object"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /contracts [post]
func (c *ContractsController) register(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	req := &api.RegisterContractRequest{}
	err := jsonutils.UnmarshalBody(request.Body, req)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	contract, err := formatters.FormatRegisterContractRequest(req)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.ucs.RegisterContract().Execute(ctx, contract)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	contract, err = c.ucs.GetContract().Execute(ctx, contract.Name, contract.Tag)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatContractResponse(contract))
}

// @Summary Set the codeHash of the given contract address
// @Description Retrieve events using hash of signature
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param address path string true "contract deployed address"
// @Param chain_id path string true "network chain id"
// @Success 200 {array} string "List of events"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /contracts/accounts/{chain_id}/{address} [post]
func (c *ContractsController) setCodeHash(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	chainID := mux.Vars(request)["chain_id"]
	address := mux.Vars(request)["address"]
	if !ethcommon.IsHexAddress(address) {
		err := errors.InvalidParameterError("expected valid address in path")
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	req := &api.SetContractCodeHashRequest{}
	err := jsonutils.UnmarshalBody(request.Body, req)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.ucs.SetContractCodeHash().Execute(ctx, chainID, ethcommon.HexToAddress(address).String(), req.CodeHash)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte("OK"))
}

// @Summary Retrieve events using hash of signature
// @Description Retrieve events using hash of signature
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param address path string true "contract deployed address"
// @Param chain_id path string true "network chain id"
// @Success 200 {object} api.GetContractEventsBySignHashResponse{} "List of events"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Events not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /contracts/accounts/{chain_id}/{address}/events [get]
func (c *ContractsController) getEvents(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	req, err := formatters.FormatGetContractEventsRequest(request)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	chainID := mux.Vars(request)["chain_id"]
	address := mux.Vars(request)["address"]
	if !ethcommon.IsHexAddress(address) {
		err := errors.InvalidParameterError("expected valid address in path")
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	abi, abiEvents, err := c.ucs.GetContractEvents().Execute(ctx, chainID, ethcommon.HexToAddress(address).String(),
		req.SigHash, req.IndexedInputCount)

	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(api.GetContractEventsBySignHashResponse{Event: abi, DefaultEvents: abiEvents})
}

// @Summary Returns a list of all tags
// @Description Returns a list of all tags from given contract name
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Success 200 {array} string "List of tags"
// @Failure 404 {object} httputil.ErrorResponse "contract not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /contracts/{name} [get]
func (c *ContractsController) getTags(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	tags, err := c.ucs.GetContractTags().Execute(ctx, mux.Vars(request)["name"])

	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	if tags == nil {
		tags = []string{}
	}
	_ = json.NewEncoder(rw).Encode(tags)
}

// @Summary Fetch registered contract data
// @Description Fetch solidity contract data by {name} and {tag}
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param name path string true "solidity contract registered name"
// @Param tag path string true "solidity contract registered tag"
// @Success 200 {object} api.ContractResponse{constructor=entities.Method,methods=[]entities.Method,events=[]entities.Event} "Contract found"
// @Failure 404 {object} httputil.ErrorResponse "Contract not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /contracts/{name}/{tag} [get]
func (c *ContractsController) getContract(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	contract, err := c.ucs.GetContract().Execute(ctx, mux.Vars(request)["name"], mux.Vars(request)["tag"])
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatContractResponse(contract))
}

// @Summary Get method signatures of registered contract
// @Description Get method signatures of registered contract by {name} and {tag}
// @Produce json
// @Security ApiKeyAuth
// @Security JWTAuth
// @Param name path string true "solidity contract registered name"
// @Param tag path string true "solidity contract registered tag"
// @Success 200 {array} string "List of signatures"
// @Failure 404 {object} httputil.ErrorResponse "Contract not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /contracts/{name}/{tag}/method-signatures [get]
func (c *ContractsController) getContractMethodSignatures(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	ctx := request.Context()

	filterMethod := request.URL.Query().Get("method")
	signatures, err := c.ucs.GetContractMethodSignatures().Execute(ctx, mux.Vars(request)["name"], mux.Vars(request)["tag"], filterMethod)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(signatures)
}
