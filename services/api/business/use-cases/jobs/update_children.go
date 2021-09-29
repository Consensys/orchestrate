package jobs

import (
	"context"
	"fmt"

	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"

	"github.com/consensys/orchestrate/pkg/errors"

	"github.com/consensys/orchestrate/pkg/types/entities"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/services/api/store"
	"github.com/consensys/orchestrate/services/api/store/models"
)

const updateChildrenComponent = "use-cases.update-children"

// createJobUseCase is a use case to create a new transaction job
type updateChildrenUseCase struct {
	db     store.DB
	logger *log.Logger
}

// NewUpdateChildrenUseCase creates a new UpdateChildrenUseCase
func NewUpdateChildrenUseCase(db store.DB) usecases.UpdateChildrenUseCase {
	return &updateChildrenUseCase{
		db:     db,
		logger: log.NewLogger().SetComponent(updateChildrenComponent),
	}
}

func (uc updateChildrenUseCase) WithDBTransaction(dbtx store.Tx) usecases.UpdateChildrenUseCase {
	uc.db = dbtx
	return &uc
}

func (uc *updateChildrenUseCase) Execute(ctx context.Context, jobUUID, parentJobUUID string, nextStatus entities.JobStatus, tenants []string) error {
	ctx = log.WithFields(ctx, log.Field("job", jobUUID), log.Field("parent_job", parentJobUUID),
		log.Field("next_status", nextStatus))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("updating sibling and/or parent jobs")

	if !entities.IsFinalJobStatus(nextStatus) {
		errMsg := "expected final job status"
		err := errors.InvalidParameterError(errMsg)
		logger.WithError(err).Error("failed to update children jobs")
		return err
	}

	jobsToUpdate, err := uc.db.Job().Search(ctx, &entities.JobFilters{
		ParentJobUUID: parentJobUUID,
		Status:        entities.StatusPending,
	}, tenants)

	if err != nil {
		return errors.FromError(err).ExtendComponent(updateChildrenComponent)
	}

	for _, jobModel := range jobsToUpdate {
		// Skip mined job which trigger the update of sibling/children
		if jobModel.UUID == jobUUID {
			continue
		}

		jobLogModel := &models.Log{
			JobID:   &jobModel.ID,
			Status:  nextStatus,
			Message: fmt.Sprintf("sibling (or parent) job %s was mined instead", jobUUID),
		}

		jobModel.Status = nextStatus
		if err := uc.db.Job().Update(ctx, jobModel); err != nil {
			return errors.FromError(err).ExtendComponent(updateChildrenComponent)
		}

		if err := uc.db.Log().Insert(ctx, jobLogModel); err != nil {
			return errors.FromError(err).ExtendComponent(updateChildrenComponent)
		}

		logger.WithField("job", jobModel.UUID).
			WithField("status", nextStatus).Debug("updated children/sibling job successfully")
	}

	logger.WithField("status", nextStatus).Info("children (and/or parent) jobs updated successfully")
	return nil
}
