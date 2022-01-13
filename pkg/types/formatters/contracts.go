package formatters

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"encoding/json"

	"github.com/consensys/orchestrate/pkg/errors"
	types "github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func FormatRegisterContractRequest(req *types.RegisterContractRequest) (*entities.Contract, error) {
	rawABI, err := json.Marshal(req.ABI)
	if err != nil {
		return nil, err
	}

	tag := req.Tag
	if tag == "" {
		tag = entities.DefaultTagValue
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(rawABI)))
	if err != nil {
		return nil, err
	}

	return &entities.Contract{
		Name:             req.Name,
		Tag:              tag,
		Bytecode:         req.Bytecode,
		DeployedBytecode: req.DeployedBytecode,
		RawABI:           string(rawABI),
		ABI:              parsedABI,
	}, nil
}

func FormatSearchContractRequest(req *http.Request) (*types.SearchContractRequest, error) {
	res := &types.SearchContractRequest{}
	var err error

	qAddress := req.URL.Query().Get("address")
	if qAddress != "" {
		addr := ethcommon.HexToAddress(qAddress)
		res.Address = &addr
	}

	qCodeHash := req.URL.Query().Get("code_hash")
	if qCodeHash != "" {
		res.CodeHash, err = hexutil.Decode(qCodeHash)
		if err != nil {
			return nil, errors.InvalidParameterError("code_hash is not hex value")
		}
	}

	if res.CodeHash == nil && res.Address == nil {
		return nil, errors.InvalidParameterError("invalid search contract request")
	}

	return res, nil
}

func FormatGetContractEventsRequest(req *http.Request) (*types.GetContractEventsRequest, error) {
	qSigHash := req.URL.Query().Get("sig_hash")
	if qSigHash == "" {
		return nil, errors.InvalidParameterError("sig_hash cannot be empty")
	}
	sigHash, err := hexutil.Decode(qSigHash)
	if err != nil {
		return nil, errors.InvalidParameterError("sig_hash is not hex value")
	}

	qIndexedInputCount := req.URL.Query().Get("indexed_input_count")
	if qIndexedInputCount == "" {
		return nil, errors.InvalidParameterError("indexed_input_count cannot be empty")
	}

	qIndexedInputCountInt, err := strconv.ParseUint(qIndexedInputCount, 10, 32)
	if err != nil {
		return nil, errors.InvalidParameterError("indexed_input_count is not valid integer")
	}

	return &types.GetContractEventsRequest{
		SigHash:           sigHash,
		IndexedInputCount: uint32(qIndexedInputCountInt),
	}, nil
}

func FormatContractResponse(contract *entities.Contract) *types.ContractResponse {
	return &types.ContractResponse{
		Name:             contract.Name,
		Tag:              contract.Tag,
		Registry:         contract.Registry,
		ABI:              contract.RawABI,
		Bytecode:         contract.Bytecode,
		DeployedBytecode: contract.DeployedBytecode,
		Constructor:      contract.Constructor,
		Methods:          contract.Methods,
		Events:           contract.Events,
	}
}
