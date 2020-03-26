package utils

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
)

const (
	defaultTag = "latest"
)

// CheckExtractArtifacts validates a request input, that is supposed to provide artifacts data
func CheckExtractArtifacts(contract *abi.Contract) (bytecode, deployedBytecode, abiBytes string, err error) {
	if contract == nil {
		return "", "", "", errors.InvalidArgError("no contract provided in request").ExtendComponent(component)
	}

	if _, err = hexutil.Decode(contract.GetBytecode()); err != nil {
		return "", "", "", errors.InvalidArgError("invalid bytecode or no contract bytecode provided in request").ExtendComponent(component)
	}

	if _, err = hexutil.Decode(contract.GetDeployedBytecode()); err != nil {
		return "", "", "", errors.InvalidArgError("invalid deployed bytecode or no contract deployed bytecode provided in requestor").ExtendComponent(component)
	}

	if contract.GetAbi() == "" {
		return "", "", "", errors.InvalidArgError("No abi provided in request").ExtendComponent(component)
	}

	compactedABI, err := contract.GetABICompacted()
	if err != nil {
		return "", "", "", errors.InvalidArgError("Failed to get compacted ABI").ExtendComponent(component)
	}

	return contract.GetBytecode(), contract.GetDeployedBytecode(), compactedABI, nil
}

// CheckExtractNameTag validates a request input, that is supposed to provide name + tag data
func CheckExtractNameTag(id *abi.ContractId) (name, tag string, err error) {
	if id == nil {
		return "", "", errors.InvalidArgError("Nil ContractId passed in the request").ExtendComponent(component)
	}

	name = id.GetName()
	if name == "" {
		return "", "", errors.InvalidArgError("No abi provided in request").ExtendComponent(component)
	}

	// Set Tag to latest if it was not set in the request
	tag = id.GetTag()
	if tag == "" {
		tag = defaultTag
	}

	return name, tag, nil
}
