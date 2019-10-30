package common

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/corestack.git/types/ethereum"
)

const (
	defaultTag = "latest"
)

// CheckExtractChainAddress validates a request input that is supposed to chain and address data
func CheckExtractChainAddress(accountInstance *common.AccountInstance) (*chain.Chain, *ethereum.Account, error) {
	if accountInstance == nil {
		return nil, nil, errors.InvalidArgError("No account instance found in request").ExtendComponent(component)
	}

	accountChain := accountInstance.GetChain()
	if accountChain == nil {
		return nil, nil, errors.InvalidArgError("No ethereum chainID found in request").ExtendComponent(component)
	}

	address := accountInstance.GetAccount()
	if address == nil {
		return nil, nil, errors.InvalidArgError("No ethereum account instance found in request").ExtendComponent(component)
	}

	return accountChain, address, nil
}

// CheckExtractArtifacts validates a request input, that is supposed to provide artifacts data
func CheckExtractArtifacts(contract *abi.Contract) (bytecode, deployedBytecode, abiBytes []byte, err error) {
	if contract == nil {
		return []byte{}, []byte{}, []byte{}, errors.InvalidArgError("No contract provided in request").ExtendComponent(component)
	}

	if contract.Bytecode == nil {
		return []byte{}, []byte{}, []byte{}, errors.InvalidArgError("No contract bytecode provided in request").ExtendComponent(component)
	}

	if contract.DeployedBytecode == nil {
		return []byte{}, []byte{}, []byte{}, errors.InvalidArgError("No contract deployed bytecode provided in request").ExtendComponent(component)
	}

	if len(contract.Abi) == 0 {
		return []byte{}, []byte{}, []byte{}, errors.InvalidArgError("No abi provided in request").ExtendComponent(component)
	}

	compactedABI, err := contract.GetABICompacted()
	if err != nil {
		return []byte{}, []byte{}, []byte{}, errors.FromError(err).ExtendComponent(component)
	}
	return contract.GetBytecode(), contract.GetDeployedBytecode(), compactedABI, nil
}

// CheckExtractNameTag validates a request input, that is supposed to provide name + tag data
func CheckExtractNameTag(id *abi.ContractId) (name, tag string, err error) {
	if id == nil {
		return "", "", errors.InvalidArgError("Nil ContractId passed in the request").ExtendComponent(component)
	}

	name = id.Name
	if name == "" {
		return "", "", errors.InvalidArgError("No abi provided in request").ExtendComponent(component)
	}

	// Set Tag to latest if it was not set in the request
	tag = id.Tag
	if tag == "" {
		tag = defaultTag
	}

	return name, tag, nil
}
