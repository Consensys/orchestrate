package contracts

import (
	"context"

	"github.com/consensys/quorum/accounts/abi"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/database"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
	"github.com/consensys/orchestrate/services/api/store/models"
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

	codeHash := crypto.Keccak256Hash(contract.DeployedBytecode)
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
		Bytecode:         contract.Bytecode.String(),
		DeployedBytecode: contract.DeployedBytecode.String(),
		Codehash:         hexutil.Encode(crypto.Keccak256(contract.DeployedBytecode)),
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

func getMethods(contractAbi *abi.ABI, deployedBytecode hexutil.Bytes, codeHash common.Hash, methodJSONs map[string]string) []*models.MethodModel {
	var methods []*models.MethodModel
	// nolint
	for _, m := range contractAbi.Methods {
		sel := sigHashToSelector(m.ID())
		if deployedBytecode != nil {
			methods = append(methods, &models.MethodModel{
				Codehash: codeHash.Hex(),
				Selector: sel,
				ABI:      methodJSONs[m.Sig()],
			})
		}
	}

	return methods
}

func getEvents(contractAbi *abi.ABI, deployedBytecode hexutil.Bytes, codeHash common.Hash, eventJSONs map[string]string) []*models.EventModel {
	var events []*models.EventModel
	// nolint
	for _, e := range contractAbi.Events {
		indexedCount := getIndexedCount(&e)
		if deployedBytecode != nil {
			events = append(events, &models.EventModel{
				Codehash:          codeHash.Hex(),
				SigHash:           e.ID().Hex(),
				IndexedInputCount: indexedCount,
				ABI:               eventJSONs[e.Sig()],
			})
		}
	}

	return events
}
