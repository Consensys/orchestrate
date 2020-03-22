package usecases

import (
	"context"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/contract-registry/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
)

const registerContractComponent = component + ".register-contract"

//go:generate mockgen -source=register_contract.go -destination=mocks/mock_register_contract.go -package=mocks

type RegisterContractUseCase interface {
	Execute(ctx context.Context, contract *abi.Contract) error
}

type PGArtifactModel struct {
	tableName struct{} `pg:"artifacts"` //nolint:unused,structcheck // reason
	models.ArtifactModel
}

// RegisterContract is a use case to register a new contract
type RegisterContract struct {
	contractDataAgent store.ContractDataAgent
}

// NewRegisterContract creates a new RegisterContract
func NewRegisterContract(contractDataAgent store.ContractDataAgent) *RegisterContract {
	return &RegisterContract{
		contractDataAgent: contractDataAgent,
	}
}

// Execute validates and registers a new contract in DB
func (usecase *RegisterContract) Execute(ctx context.Context, contract *abi.Contract) error {
	bytecode, deployedBytecode, abiRaw, err := utils.CheckExtractArtifacts(contract)
	if err != nil {
		return errors.FromError(err).ExtendComponent(registerContractComponent)
	}

	name, tagName, err := utils.CheckExtractNameTag(contract.Id)
	if err != nil {
		return errors.FromError(err).ExtendComponent(registerContractComponent)
	}

	d, err := hexutil.Decode(deployedBytecode)
	if err != nil {
		return errors.InvalidArgError("Could not decode deployedBytecode").ExtendComponent(registerContractComponent)
	}

	codeHash := crypto.Keccak256Hash(d)
	contractAbi, err := contract.ToABI()
	if err != nil {
		return errors.FromError(err).ExtendComponent(registerContractComponent)
	}
	methodJSONs, eventJSONs, err := utils.ParseJSONABI(abiRaw)
	if err != nil {
		return errors.FromError(err).ExtendComponent(registerContractComponent)
	}

	methods := getMethods(contractAbi, deployedBytecode, codeHash, methodJSONs)
	events := getEvents(contractAbi, deployedBytecode, codeHash, eventJSONs)

	err = usecase.contractDataAgent.Insert(ctx, name, tagName, abiRaw, bytecode, deployedBytecode, hexutil.Encode(crypto.Keccak256(d)), &methods, &events)
	if err != nil {
		return errors.FromError(err).ExtendComponent(registerContractComponent)
	}

	return nil
}

func getMethods(contractAbi *ethabi.ABI, deployedBytecode string, codeHash common.Hash, methodJSONs map[string]string) []*models.MethodModel {
	var methods []*models.MethodModel
	for _, m := range contractAbi.Methods {
		// Register methods for this bytecode
		method := m
		sel := utils.SigHashToSelector(method.ID())
		if deployedBytecode != "" {
			methods = append(methods, &models.MethodModel{
				Codehash: codeHash.Hex(),
				Selector: sel,
				ABI:      methodJSONs[method.Sig()],
			})
		}
	}

	return methods
}

func getEvents(contractAbi *ethabi.ABI, deployedBytecode string, codeHash common.Hash, eventJSONs map[string]string) []*models.EventModel {
	var events []*models.EventModel
	for _, e := range contractAbi.Events {
		event := e
		indexedCount := utils.GetIndexedCount(event)
		// Register events for this bytecode
		if deployedBytecode != "" {
			events = append(events, &models.EventModel{
				Codehash:          codeHash.Hex(),
				SigHash:           event.ID().Hex(),
				IndexedInputCount: indexedCount,
				ABI:               eventJSONs[event.Sig()],
			})
		}
	}

	return events
}
