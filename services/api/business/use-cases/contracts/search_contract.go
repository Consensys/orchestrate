package contracts

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/entities"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const searchContractComponent = "use-cases.search-contract"

type searchContractUseCase struct {
	agent  store.ContractAgent
	logger *log.Logger
}

func NewSearchContractUseCase(agent store.ContractAgent) usecases.SearchContractUseCase {
	return &searchContractUseCase{
		agent:  agent,
		logger: log.NewLogger().SetComponent(searchContractComponent),
	}
}

func (uc *searchContractUseCase) Execute(ctx context.Context, codehash hexutil.Bytes, address *ethcommon.Address) (*entities.Contract, error) {
	logger := uc.logger.WithContext(ctx)

	var contract *entities.Contract
	var err error
	switch {
	case address != nil:
		contract, err = uc.agent.FindOneByAddress(ctx, address.String())
	case codehash != nil:
		contract, err = uc.agent.FindOneByCodeHash(ctx, codehash.String())
	}

	if err != nil {
		uc.logger.WithError(err).Error("no contract found")
		return nil, errors.FromError(err).ExtendComponent(searchContractComponent)
	}

	logger.WithField("name", contract.Name).WithField("tag", contract.Tag).Debug("contract found")
	return contract, nil
}
