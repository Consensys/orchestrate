package contracts

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const getEventsComponent = "use-cases.get-events"

type getEventsUseCase struct {
	agent store.EventAgent
}

func NewGetEventsUseCase(agent store.EventAgent) usecases.GetContractEventsUseCase {
	return &getEventsUseCase{
		agent: agent,
	}
}

// Execute validates and registers a new contract in DB
func (usecase *getEventsUseCase) Execute(ctx context.Context, chainID, address, sighash string, indexedInputCount uint32) (abi string, eventsABI []string, err error) {
	logger := log.WithContext(ctx).WithField("chainID", chainID).WithField("address", address).
		WithField("sig_hash", utils.ShortString(sighash, 10))
	logger.Debug("get events starting...")

	eventModel, err := usecase.agent.FindOneByAccountAndSigHash(ctx, chainID, address, sighash, indexedInputCount)
	if err != nil && !errors.IsNotFoundError(err) {
		return "", nil, errors.FromError(err).ExtendComponent(getEventsComponent)
	}

	if eventModel != nil {
		logger.Debug("get events executed successfully")
		return eventModel.ABI, nil, nil
	}

	defaultEventModels, err := usecase.agent.FindDefaultBySigHash(ctx, sighash, indexedInputCount)
	if err != nil {
		return "", nil, errors.FromError(err).ExtendComponent(getEventsComponent)
	}

	for _, e := range defaultEventModels {
		eventsABI = append(eventsABI, e.ABI)
	}

	logger.Debug("get events executed successfully")
	return "", eventsABI, nil
}
