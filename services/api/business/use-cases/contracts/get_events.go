package contracts

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/ConsenSys/orchestrate/services/api/store"
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
func (uc *getEventsUseCase) Execute(ctx context.Context, chainID, address, sighash string, indexedInputCount uint32) (abi string, eventsABI []string, err error) {
	ctx = log.WithFields(ctx, log.Field("chain_id", chainID), log.Field("address", address))
	logger := uc.logger.WithContext(ctx)

	eventModel, err := uc.agent.FindOneByAccountAndSigHash(ctx, chainID, address, sighash, indexedInputCount)
	if err != nil && !errors.IsNotFoundError(err) {
		return "", nil, errors.FromError(err).ExtendComponent(getEventsComponent)
	}

	if eventModel != nil {
		logger.Debug("events were fetched successfully")
		return eventModel.ABI, nil, nil
	}

	defaultEventModels, err := uc.agent.FindDefaultBySigHash(ctx, sighash, indexedInputCount)
	if err != nil {
		return "", nil, errors.FromError(err).ExtendComponent(getEventsComponent)
	}

	for _, e := range defaultEventModels {
		eventsABI = append(eventsABI, e.ABI)
	}

	logger.Debug("default events were fetched successfully")
	return "", eventsABI, nil
}
