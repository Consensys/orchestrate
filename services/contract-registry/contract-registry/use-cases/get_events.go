package usecases

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store"
)

const getEventsComponent = component + ".get-events"

//go:generate mockgen -source=get_events.go -destination=mocks/mock_get_events.go -package=mocks

type GetEventsUseCase interface {
	Execute(ctx context.Context, account *common.AccountInstance, sighash string, indexedInputCount uint32) (abi string, eventsABI []string, err error)
}

// GetEvents is a use case to get events
type GetEvents struct {
	eventDataAgent store.EventDataAgent
}

// NewGetEvents creates a new GetEvents
func NewGetEvents(eventDataAgent store.EventDataAgent) *GetEvents {
	return &GetEvents{
		eventDataAgent: eventDataAgent,
	}
}

// Execute validates and registers a new contract in DB
func (usecase *GetEvents) Execute(ctx context.Context, account *common.AccountInstance, sighash string, indexedInputCount uint32) (abi string, eventsABI []string, err error) {
	event, err := usecase.eventDataAgent.FindOneByAccountAndSigHash(ctx, account, sighash, indexedInputCount)
	if errors.IsConnectionError(err) {
		return "", nil, errors.FromError(err).ExtendComponent(getEventsComponent)
	}
	if event != nil {
		return event.ABI, nil, nil
	}

	defaultEvents, err := usecase.eventDataAgent.FindDefaultBySigHash(ctx, sighash, indexedInputCount)
	if err != nil {
		return "", nil, errors.FromError(err).ExtendComponent(getEventsComponent)
	}

	for _, e := range defaultEvents {
		eventsABI = append(eventsABI, e.ABI)
	}

	return "", eventsABI, nil
}
