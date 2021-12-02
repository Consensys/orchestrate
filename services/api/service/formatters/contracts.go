package formatters

import (
	"net/http"
	"strconv"

	"github.com/consensys/orchestrate/pkg/encoding/json"
	"github.com/consensys/orchestrate/pkg/errors"
	types "github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
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

	return &entities.Contract{
		Name:             req.Name,
		Tag:              tag,
		Bytecode:         req.Bytecode,
		DeployedBytecode: req.DeployedBytecode,
		ABI:              string(rawABI),
	}, nil
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
		ABI:              contract.ABI,
		Bytecode:         contract.Bytecode,
		DeployedBytecode: contract.DeployedBytecode,
		Constructor:      contract.Constructor,
		Methods:          contract.Methods,
		Events:           contract.Events,
	}
}
