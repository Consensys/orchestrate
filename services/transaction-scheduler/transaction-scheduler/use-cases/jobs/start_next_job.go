package jobs

import (
	"context"
	"strconv"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
)

//go:generate mockgen -source=start_next_job.go -destination=mocks/start_next_job.go -package=mocks

const startNextJobComponent = "use-cases.next-job-start"

// startNextJobUseCase is a use case to get a job
type startNextJobUseCase struct {
	db              store.DB
	startJobUseCase usecases.StartJobUseCase
}

// NewStartNextJobUseCase creates a new StartNextJobUseCase
func NewStartNextJobUseCase(db store.DB, startJobUC usecases.StartJobUseCase) usecases.StartNextJobUseCase {
	return &startNextJobUseCase{
		db:              db,
		startJobUseCase: startJobUC,
	}
}

// Execute gets a job
func (uc *startNextJobUseCase) Execute(ctx context.Context, jobUUID string, tenants []string) error {
	jobModel, err := uc.db.Job().FindOneByUUID(ctx, jobUUID, tenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(startNextJobComponent)
	}

	if jobModel.NextJobUUID == "" {
		return errors.DataError("job %s does not have a next job to start", jobModel.NextJobUUID)
	}

	logger := log.WithContext(ctx).
		WithField("job_uuid", jobUUID).
		WithField("next_job_uuid", jobModel.NextJobUUID)

	logger.Debug("start next job use-case")

	nextJobModel, err := uc.db.Job().FindOneByUUID(ctx, jobModel.NextJobUUID, tenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(startNextJobComponent)
	}

	switch nextJobModel.Type {
	case utils.OrionMarkingTransaction:
		err = uc.handleOrionMarkingTx(ctx, jobModel, nextJobModel)
	case utils.TesseraMarkingTransaction:
		err = uc.handleTesseraMarkingTx(ctx, jobModel, nextJobModel)
	}

	if err != nil {
		return errors.FromError(err).ExtendComponent(startNextJobComponent)
	}

	return uc.startJobUseCase.Execute(ctx, nextJobModel.UUID, tenants)
}

func (uc *startNextJobUseCase) handleOrionMarkingTx(ctx context.Context, prevJobModel, jobModel *models.Job) error {
	if prevJobModel.Type != utils.OrionEEATransaction {
		return errors.DataError("expected previous job as type: %s", utils.OrionEEATransaction)
	}

	prevJobEntity := parsers.NewJobEntityFromModels(prevJobModel)
	if prevJobEntity.GetStatus() != utils.StatusStored {
		return errors.DataError("expected previous job status as: STORED")
	}

	jobModel.Transaction.Data = prevJobModel.Transaction.Hash
	return uc.db.Transaction().Update(ctx, jobModel.Transaction)
}

func (uc *startNextJobUseCase) handleTesseraMarkingTx(ctx context.Context, prevJobModel, jobModel *models.Job) error {
	if prevJobModel.Type != utils.TesseraPrivateTransaction {
		return errors.DataError("expected previous job as type: %s", utils.TesseraPrivateTransaction)
	}

	prevJobEntity := parsers.NewJobEntityFromModels(prevJobModel)
	if prevJobEntity.GetStatus() != utils.StatusStored {
		return errors.DataError("expected previous job status as: STORED")
	}

	jobModel.Transaction.Data = prevJobModel.Transaction.EnclaveKey
	gas, err := strconv.ParseInt(prevJobModel.Transaction.Gas, 10, 64)
	if err == nil && gas < utils.TesseraGasLimit {
		jobModel.Transaction.Gas = strconv.Itoa(utils.TesseraGasLimit)
	} else {
		jobModel.Transaction.Gas = prevJobModel.Transaction.Gas
	}

	return uc.db.Transaction().Update(ctx, jobModel.Transaction)
}
