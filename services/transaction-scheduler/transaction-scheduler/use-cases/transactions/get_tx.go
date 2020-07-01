package transactions

import (
	"context"
	"sort"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
)

//go:generate mockgen -source=get_tx.go -destination=mocks/get_tx.go -package=mocks

const getTxComponent = "use-cases.get-tx"

type GetTxUseCase interface {
	Execute(ctx context.Context, txRequestUUID string, tenants []string) (*entities.TxRequest, error)
}

// getTxUseCase is a use case to get a transaction request
type getTxUseCase struct {
	db                 store.DB
	getScheduleUsecase schedules.GetScheduleUseCase
}

// NewGetTxUseCase creates a new GetTxUseCase
func NewGetTxUseCase(db store.DB, getScheduleUsecase schedules.GetScheduleUseCase) GetTxUseCase {
	return &getTxUseCase{
		db:                 db,
		getScheduleUsecase: getScheduleUsecase,
	}
}

// Execute gets a transaction request
func (uc *getTxUseCase) Execute(ctx context.Context, txRequestUUID string, tenants []string) (*entities.TxRequest, error) {
	logger := log.WithContext(ctx).WithField("tx_request_uuid", txRequestUUID)
	logger.Debug("getting transaction request")

	txRequestModel, err := uc.db.TransactionRequest().FindOneByUUID(ctx, txRequestUUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getTxComponent)
	}

	// Sort schedules by timestamp as we want to retrieve the last one, here it's sorted in decreasing time
	// So txRequestModel.Schedules[0] is the latest
	sort.SliceStable(txRequestModel.Schedules, func(i, j int) bool {
		return txRequestModel.Schedules[i].CreatedAt.After(txRequestModel.Schedules[j].CreatedAt)
	})

	txRequest := &entities.TxRequest{
		UUID:           txRequestModel.UUID,
		IdempotencyKey: txRequestModel.IdempotencyKey,
		Params:         txRequestModel.Params,
		CreatedAt:      txRequestModel.CreatedAt,
	}
	txRequest.Schedule, err = uc.getScheduleUsecase.Execute(ctx, txRequestModel.Schedules[0].UUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getTxComponent)
	}

	logger.Info("transaction request found successfully")

	return txRequest, nil
}
