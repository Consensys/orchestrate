package contracts

import (
	"context"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store/models"
)

const registerContractComponent = "use-cases.register-contract"

type registerContractUseCase struct {
	db store.DB
}

func NewRegisterContractUseCase(agent store.DB) usecases.RegisterContractUseCase {
	return &registerContractUseCase{
		db: agent,
	}
}

func (uc *registerContractUseCase) Execute(ctx context.Context, contract *entities.Contract) error {
	logger := log.WithContext(ctx).WithField("contract", contract.ID.Short())
	logger.Debug("registering contract starting...")

	abiRaw, err := contract.GetABICompacted()
	if err != nil {
		return errors.FromError(err).ExtendComponent(registerContractComponent)
	}

	dByteCode, err := hexutil.Decode(contract.DeployedBytecode)
	if err != nil {
		logger.WithError(err).Error("cannot decode deployedByteCode")
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
		logger.WithError(err).Error("cannot parse json ABI")
		return errors.FromError(err).ExtendComponent(registerContractComponent)
	}

	repository := &models.RepositoryModel{
		Name: contract.ID.Name,
	}
	artifact := &models.ArtifactModel{
		ABI:              abiRaw,
		Bytecode:         contract.Bytecode,
		DeployedBytecode: contract.DeployedBytecode,
		Codehash:         hexutil.Encode(crypto.Keccak256(dByteCode)),
	}

	err = database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
		dbtx := tx.(store.Tx)
		if der := selectOrInsertRepository(ctx, dbtx, repository); der != nil {
			return der
		}
		if der := selectOrInsertArtifact(ctx, dbtx, artifact); der != nil {
			return der
		}

		tag := &models.TagModel{
			Name:         contract.ID.Tag,
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

func selectOrInsertRepository(ctx context.Context, dbtx store.Tx, repository *models.RepositoryModel) error {
	repo, err := dbtx.Repository().FindOneAndLock(ctx, repository.Name)
	if err != nil && !errors.IsNotFoundError(err) {
		return err
	}

	if repo != nil {
		repository.ID = repo.ID
		return nil
	}

	return dbtx.Repository().Insert(ctx, repository)
}

func selectOrInsertArtifact(ctx context.Context, dbtx store.Tx, artifact *models.ArtifactModel) error {
	arti, err := dbtx.Artifact().FindOneByABIAndCodeHash(ctx, artifact.ABI, artifact.Codehash)
	if err != nil && !errors.IsNotFoundError(err) {
		return err
	}

	if arti != nil {
		artifact.ID = arti.ID
		return nil
	}

	return dbtx.Artifact().Insert(ctx, artifact)
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
