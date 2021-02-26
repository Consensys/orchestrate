package transactions

import (
	"context"

	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	"github.com/ConsenSys/orchestrate/services/api/store"
)

const getTxComponent = "use-cases.get-tx"

// getTxUseCase is a use case to get a transaction request
type getTxUseCase struct {
	db                 store.DB
	getScheduleUsecase usecases.GetScheduleUseCase
	logger             *log.Logger
}

// NewGetTxUseCase creates a new GetTxUseCase
func NewGetTxUseCase(db store.DB, getScheduleUsecase usecases.GetScheduleUseCase) usecases.GetTxUseCase {
	return &getTxUseCase{
		db:                 db,
		getScheduleUsecase: getScheduleUsecase,
		logger:             log.NewLogger().SetComponent(getTxComponent),
	}
}

// Execute gets a transaction request
func (uc *getTxUseCase) Execute(ctx context.Context, scheduleUUID string, tenants []string) (*entities.TxRequest, error) {
	ctx = log.WithFields(ctx, log.Field("schedule", scheduleUUID))

	txRequestModel, err := uc.db.TransactionRequest().FindOneByUUID(ctx, scheduleUUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getTxComponent)
	}

	txRequest := &entities.TxRequest{
		IdempotencyKey: txRequestModel.IdempotencyKey,
		ChainName:      txRequestModel.ChainName,
		Params:         txRequestModel.Params,
		CreatedAt:      txRequestModel.CreatedAt,
	}
	txRequest.Schedule, err = uc.getScheduleUsecase.Execute(ctx, txRequestModel.Schedule.UUID, tenants)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getTxComponent)
	}

	uc.logger.WithContext(ctx).Debug("transaction request found successfully")
	return txRequest, nil
}
