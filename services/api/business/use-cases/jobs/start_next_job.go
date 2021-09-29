package jobs

import (
	"context"
	"strconv"

	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store/models"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/services/api/store"
)

const startNextJobComponent = "use-cases.next-job-start"

type startNextJobUseCase struct {
	db              store.DB
	startJobUseCase usecases.StartJobUseCase
	logger          *log.Logger
}

func NewStartNextJobUseCase(db store.DB, startJobUC usecases.StartJobUseCase) usecases.StartNextJobUseCase {
	return &startNextJobUseCase{
		db:              db,
		startJobUseCase: startJobUC,
		logger:          log.NewLogger().SetComponent(startNextJobComponent),
	}
}

// Execute gets a job
func (uc *startNextJobUseCase) Execute(ctx context.Context, jobUUID string, tenants []string) error {
	ctx = log.WithFields(ctx, log.Field("job", jobUUID))
	jobModel, err := uc.db.Job().FindOneByUUID(ctx, jobUUID, tenants, false)
	if err != nil {
		return errors.FromError(err).ExtendComponent(startNextJobComponent)
	}

	if jobModel.NextJobUUID == "" {
		return errors.DataError("job %s does not have a next job to start", jobModel.NextJobUUID)
	}

	logger := uc.logger.WithContext(ctx).WithField("next_job", jobModel.NextJobUUID)
	logger.Debug("start next job use-case")

	nextJobModel, err := uc.db.Job().FindOneByUUID(ctx, jobModel.NextJobUUID, tenants, false)
	if err != nil {
		return errors.FromError(err).ExtendComponent(startNextJobComponent)
	}

	switch nextJobModel.Type {
	case entities.OrionMarkingTransaction:
		err = uc.handleOrionMarkingTx(ctx, jobModel, nextJobModel)
	case entities.TesseraMarkingTransaction:
		err = uc.handleTesseraMarkingTx(ctx, jobModel, nextJobModel)
	}

	if err != nil {
		logger.WithError(err).Error("failed to validate next transaction data")
		return errors.FromError(err).ExtendComponent(startNextJobComponent)
	}

	return uc.startJobUseCase.Execute(ctx, nextJobModel.UUID, tenants)
}

func (uc *startNextJobUseCase) handleOrionMarkingTx(ctx context.Context, prevJobModel, jobModel *models.Job) error {
	if prevJobModel.Type != entities.OrionEEATransaction {
		return errors.DataError("expected previous job as type: %s", entities.OrionEEATransaction)
	}

	prevJobEntity := parsers.NewJobEntityFromModels(prevJobModel)
	if prevJobEntity.Status != entities.StatusStored {
		return errors.DataError("expected previous job status as: STORED")
	}

	jobModel.Transaction.Data = prevJobModel.Transaction.Hash
	return uc.db.Transaction().Update(ctx, jobModel.Transaction)
}

func (uc *startNextJobUseCase) handleTesseraMarkingTx(ctx context.Context, prevJobModel, jobModel *models.Job) error {
	if prevJobModel.Type != entities.TesseraPrivateTransaction {
		return errors.DataError("expected previous job as type: %s", entities.TesseraPrivateTransaction)
	}

	prevJobEntity := parsers.NewJobEntityFromModels(prevJobModel)
	if prevJobEntity.Status != entities.StatusStored {
		return errors.DataError("expected previous job status as: STORED")
	}

	jobModel.Transaction.Data = prevJobModel.Transaction.EnclaveKey
	gas, err := strconv.ParseInt(prevJobModel.Transaction.Gas, 10, 64)
	if err == nil && gas < entities.TesseraGasLimit {
		jobModel.Transaction.Gas = strconv.Itoa(entities.TesseraGasLimit)
	} else {
		jobModel.Transaction.Gas = prevJobModel.Transaction.Gas
	}

	return uc.db.Transaction().Update(ctx, jobModel.Transaction)
}
