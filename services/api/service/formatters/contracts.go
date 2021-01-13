package formatters

import (
	"net/http"
	"strconv"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
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
	if !utils.IsHexString(qSigHash) {
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
		SigHash:           qSigHash,
		IndexedInputCount: uint32(qIndexedInputCountInt),
	}, nil
}
