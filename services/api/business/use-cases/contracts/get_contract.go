package contracts

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/consensys/orchestrate/services/api/store"
	"github.com/consensys/quorum/common/hexutil"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/entities"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
)

const getContractComponent = "use-cases.get-contract"

type getContractUseCase struct {
	agent  store.ArtifactAgent
	logger *log.Logger
}

func NewGetContractUseCase(agent store.ArtifactAgent) usecases.GetContractUseCase {
	return &getContractUseCase{
		agent:  agent,
		logger: log.NewLogger().SetComponent(getContractComponent),
	}
}

// Execute gets a contract from DB
func (uc *getContractUseCase) Execute(ctx context.Context, name, tag string) (*entities.Contract, error) {
	ctx = log.WithFields(ctx, log.Field("contract_name", name), log.Field("contract_tag", name))
	logger := uc.logger.WithContext(ctx)

	artifact, err := uc.agent.FindOneByNameAndTag(ctx, name, tag)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getContractComponent)
	}

	parsedABI, err := abi.JSON(strings.NewReader(artifact.ABI))
	if err != nil {
		errMessage := "failed to parse contract abi"
		uc.logger.WithError(err).Error(errMessage)
		return nil, errors.DataCorruptedError(errMessage).ExtendComponent(getContractComponent)
	}

	contract := &entities.Contract{
		Name:             name,
		Tag:              tag,
		RawABI:           artifact.ABI,
		ABI:              parsedABI,
		Bytecode:         hexutil.MustDecode(artifact.Bytecode),
		DeployedBytecode: hexutil.MustDecode(artifact.DeployedBytecode),
	}

	// nolint
	for _, method := range parsedABI.Methods {
		contract.Methods = append(contract.Methods, entities.ABIComponent{
			Signature: method.Sig,
		})
	}

	// nolint
	for _, event := range parsedABI.Events {
		contract.Events = append(contract.Events, entities.ABIComponent{
			Signature: event.Sig,
		})
	}

	contract.Constructor = entities.ABIComponent{Signature: parsedABI.Constructor.Sig}

	logger.Debug("contract found successfully")
	return contract, nil
}
