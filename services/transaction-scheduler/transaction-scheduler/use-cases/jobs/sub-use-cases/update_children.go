package subusecases

import (
	"context"
	"fmt"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
)

//go:generate mockgen -source=update_children.go -destination=mocks/update_children.go -package=mocks

const updateChildrenComponent = "use-cases.update-children"

type UpdateChildrenUseCase interface {
	Execute(ctx context.Context, jobUUID, parentJobUUID, nextStatus string, tenants []string) error
	WithDBTransaction(dbtx store.Tx) UpdateChildrenUseCase
}

// createJobUseCase is a use case to create a new transaction job
type updateChildrenUseCase struct {
	db store.DB
}

// NewUpdateChildrenUseCase creates a new UpdateChildrenUseCase
func NewUpdateChildrenUseCase(db store.DB) UpdateChildrenUseCase {
	return &updateChildrenUseCase{
		db: db,
	}
}

func (uc updateChildrenUseCase) WithDBTransaction(dbtx store.Tx) UpdateChildrenUseCase {
	uc.db = dbtx
	return &uc
}

// Execute updates all children of a job to NEVER_MINED
func (uc *updateChildrenUseCase) Execute(ctx context.Context, jobUUID, parentJobUUID, nextStatus string, tenants []string) error {
	logger := log.WithContext(ctx).WithField("job_uuid", jobUUID).WithField("tenants", tenants)
	logger.Debug("updating sibling and/or parent jobs")

	if parentJobUUID == "" {
		parentJobUUID = jobUUID
	}

	err := uc.db.Job().LockOneByUUID(ctx, parentJobUUID)
	if err != nil {
		return errors.FromError(err).ExtendComponent(updateChildrenComponent)
	}

	jobsToUpdate, err := uc.db.Job().Search(ctx, &entities.JobFilters{ParentJobUUID: parentJobUUID}, tenants)
	if err != nil {
		return errors.FromError(err).ExtendComponent(updateChildrenComponent)
	}

	for _, jobModel := range jobsToUpdate {
		status := parsers.NewJobEntityFromModels(jobModel).GetStatus()
		if jobModel.UUID != jobUUID && status == utils.StatusPending {
			jobLogModel := &models.Log{
				JobID:   &jobModel.ID,
				Status:  nextStatus,
				Message: fmt.Sprintf("sibling (or parent) job %s was mined instead", jobUUID),
			}
			err := uc.db.Log().Insert(ctx, jobLogModel)
			if err != nil {
				return errors.FromError(err).ExtendComponent(updateChildrenComponent)
			}
		}
	}

	logger.Info("children (and/or parent) jobs updated successfully")
	return nil
}
