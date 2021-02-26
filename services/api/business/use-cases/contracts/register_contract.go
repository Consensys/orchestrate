package contracts

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/database"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/api/business/parsers"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/ConsenSys/orchestrate/services/api/store"
	"github.com/ConsenSys/orchestrate/services/api/store/models"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const registerContractComponent = "use-cases.register-contract"

type registerContractUseCase struct {
	db     store.DB
	logger *log.Logger
}

func NewRegisterContractUseCase(agent store.DB) usecases.RegisterContractUseCase {
	return &registerContractUseCase{
		db:     agent,
		logger: log.NewLogger().SetComponent(registerContractComponent),
	}
}

func (uc *registerContractUseCase) Execute(ctx context.Context, contract *entities.Contract) error {
	ctx = log.WithFields(ctx, log.Field("contract_id", contract.Short()))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("registering contract starting...")

	abiRaw, err := contract.GetABICompacted()
	if err != nil {
		return errors.FromError(err).ExtendComponent(registerContractComponent)
	}

	dByteCode, err := hexutil.Decode(contract.DeployedBytecode)
	if err != nil {
		logger.WithError(err).Error("failed to decode deployedByteCode")
		return errors.FromError(err).ExtendComponent(registerContractComponent)
	}

	codeHash := crypto.Keccak256Hash(dByteCode)
	contractAbi, err := contract.ToABI()
	if err != nil {
		logger.WithError(err).Error("invalid ABI value")
		return errors.FromError(err).ExtendComponent(registerContractComponent)
	}

	methodJSONs, eventJSONs, err := parsers.ParseJSONABI(abiRaw)
	if err != nil {
		logger.WithError(err).Error("failed to parse json ABI")
		return errors.FromError(err).ExtendComponent(registerContractComponent)
	}

	repository := &models.RepositoryModel{
		Name: contract.Name,
	}
	artifact := &models.ArtifactModel{
		ABI:              abiRaw,
		Bytecode:         contract.Bytecode,
		DeployedBytecode: contract.DeployedBytecode,
		Codehash:         hexutil.Encode(crypto.Keccak256(dByteCode)),
	}

	err = database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
		// @TODO Improve duplicate inserts when `DeployedBytecode` and `Name` and `Tag` already exists
		dbtx := tx.(store.Tx)
		if der := dbtx.Repository().SelectOrInsert(ctx, repository); der != nil {
			return der
		}
		if der := dbtx.Artifact().SelectOrInsert(ctx, artifact); der != nil {
			return der
		}

		tag := &models.TagModel{
			Name:         contract.Tag,
			RepositoryID: repository.ID,
			ArtifactID:   artifact.ID,
		}

		if der := dbtx.Tag().Insert(ctx, tag); der != nil {
			return der
		}

		methods := getMethods(contractAbi, contract.DeployedBytecode, codeHash, methodJSONs)
		if len(methods) > 0 {
			if der := dbtx.Method().InsertMultiple(ctx, methods); der != nil {
				return der
			}
		}

		events := getEvents(contractAbi, contract.DeployedBytecode, codeHash, eventJSONs)
		if len(events) > 0 {
			if der := dbtx.Event().InsertMultiple(ctx, events); der != nil {
				return der
			}
		}

		return nil
	})

	if err != nil {
		return errors.FromError(err).ExtendComponent(registerContractComponent)
	}

	logger.Info("contract registered successfully")
	return nil
}

func getMethods(contractAbi *ethabi.ABI, deployedBytecode string, codeHash common.Hash, methodJSONs map[string]string) []*models.MethodModel {
	var methods []*models.MethodModel
	for _, m := range contractAbi.Methods {
		method := m
		sel := sigHashToSelector(method.ID())
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
		indexedCount := getIndexedCount(event)
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
