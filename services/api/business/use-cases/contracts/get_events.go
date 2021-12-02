package contracts

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const getEventsComponent = "use-cases.get-events"

type getEventsUseCase struct {
	agent  store.EventAgent
	logger *log.Logger
}

func NewGetEventsUseCase(agent store.EventAgent) usecases.GetContractEventsUseCase {
	return &getEventsUseCase{
		agent:  agent,
		logger: log.NewLogger().SetComponent(getEventsComponent),
	}
}

// Execute validates and registers a new contract in DB
func (uc *getEventsUseCase) Execute(ctx context.Context, chainID string, address ethcommon.Address, sighash hexutil.Bytes, indexedInputCount uint32) (abi string, eventsABI []string, err error) {
	ctx = log.WithFields(ctx, log.Field("chain_id", chainID), log.Field("address", address))
	logger := uc.logger.WithContext(ctx)

	eventModel, err := uc.agent.FindOneByAccountAndSigHash(ctx, chainID, address.Hex(), sighash.String(), indexedInputCount)
	if err != nil && !errors.IsNotFoundError(err) {
		return "", nil, errors.FromError(err).ExtendComponent(getEventsComponent)
	}

	if eventModel != nil {
		logger.Debug("events were fetched successfully")
		return eventModel.ABI, nil, nil
	}

	defaultEventModels, err := uc.agent.FindDefaultBySigHash(ctx, sighash.String(), indexedInputCount)
	if err != nil {
		return "", nil, errors.FromError(err).ExtendComponent(getEventsComponent)
	}

	for _, e := range defaultEventModels {
		eventsABI = append(eventsABI, e.ABI)
	}

	logger.Debug("default events were fetched successfully")
	return "", eventsABI, nil
}
