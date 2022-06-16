package transactions

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
)

type speedUpTxUseCase struct {
	getTxUC      usecases.GetTxUseCase
	retryJobTxUC usecases.RetryJobTxUseCase
	logger       *log.Logger
}

func NewSpeedUpTxUseCase(getTxUC usecases.GetTxUseCase, retryJobTxUC usecases.RetryJobTxUseCase) usecases.SpeedUpTxUseCase {
	return &speedUpTxUseCase{
		getTxUC:      getTxUC,
		retryJobTxUC: retryJobTxUC,
		logger:       log.NewLogger().SetComponent("use-cases.speed-up-tx"),
	}
}

func (uc *speedUpTxUseCase) Execute(ctx context.Context, scheduleUUID string, gasIncrement float64, userInfo *multitenancy.UserInfo) (*entities.TxRequest, error) {
	ctx = log.WithFields(
		ctx,
		log.Field("schedule", scheduleUUID),
	)
	logger := uc.logger.WithContext(ctx)
	logger.Debug("speeding up transaction")

	tx, err := uc.getTxUC.Execute(ctx, scheduleUUID, userInfo)
	if err != nil {
		return nil, err
	}

	if tx.Params.Protocol != "" {
		errMsg := "speed up is not supported for private transactions"
		logger.Error(errMsg)
		return nil, errors.InvalidParameterError(errMsg)
	}

	if tx.InternalData != nil && tx.InternalData.OneTimeKey {
		errMsg := "speed up is not supported for oneTimeKey transactions"
		logger.Error(errMsg)
		return nil, errors.InvalidParameterError(errMsg)
	}

	job := tx.Schedule.Jobs[0]
	err = uc.retryJobTxUC.Execute(ctx, job.UUID, gasIncrement, job.Transaction.Data, userInfo)
	if err != nil {
		return nil, err
	}

	txRequest, err := uc.getTxUC.Execute(ctx, scheduleUUID, userInfo)
	if err != nil {
		return nil, err
	}

	logger.WithField("schedule", txRequest.Schedule.UUID).Info("speed-up transaction was sent successfully")
	return txRequest, nil
}
