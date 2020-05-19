package jobs

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	tsorm "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/orm"
)

//go:generate mockgen -source=update_job.go -destination=mocks/update_job.go -package=mocks

const updateJobComponent = "use-cases.update-job"

type UpdateJobUseCase interface {
	Execute(ctx context.Context, jobUUID string, jobRequest *types.JobUpdateRequest, tenantID string) (*types.JobResponse, error)
}

// updateJobUseCase is a use case to create a new transaction job
type updateJobUseCase struct {
	db  store.DB
	orm tsorm.ORM
}

// NewUpdateJobUseCase creates a new UpdateJobUseCase
func NewUpdateJobUseCase(db store.DB, orm tsorm.ORM) UpdateJobUseCase {
	return &updateJobUseCase{
		db:  db,
		orm: orm,
	}
}

// Execute validates and creates a new transaction job
func (uc *updateJobUseCase) Execute(ctx context.Context, jobUUID string, request *types.JobUpdateRequest, tenantID string) (*types.JobResponse, error) {
	log.WithContext(ctx).
		WithField("tenant_id", tenantID).
		WithField("job_uuid", jobUUID).
		Debug("update job")

	err := utils.GetValidator().Struct(request)
	if err != nil {
		errMessage := "failed to validate update job request"
		log.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage).ExtendComponent(updateJobComponent)
	}

	job, err := uc.db.Job().FindOneByUUID(ctx, jobUUID, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(updateJobComponent)
	}

	updateJobFromRequest(job, request)
	updateTxFromRequest(job.Transaction, request)

	job.Logs = append(job.Logs, &models.Log{
		Status: job.GetStatus(),
		// @TODO Improve entry log message
		Message:   fmt.Sprintf("updated job %q", request),
		CreatedAt: time.Now(),
	})

	job.TransactionID = nil // Indicates that Tx should be updated
	if err := uc.orm.InsertOrUpdateJob(ctx, uc.db, job); err != nil {
		return nil, errors.FromError(err).ExtendComponent(createJobComponent)
	}

	log.WithContext(ctx).
		WithField("job_uuid", job.UUID).
		Info("job updated successfully")

	return &types.JobResponse{
		UUID:        job.UUID,
		Transaction: request.Transaction,
		Status:      job.GetStatus(),
		CreatedAt:   job.CreatedAt,
	}, nil
}

func updateJobFromRequest(job *models.Job, request *types.JobUpdateRequest) {
	job.Labels = request.Labels
}

// @TODO Improve next by an smarted solution
func updateTxFromRequest(tx *models.Transaction, request *types.JobUpdateRequest) {
	if request.Transaction.Hash != "" {
		tx.Hash = request.Transaction.Hash
	}
	if request.Transaction.From != "" {
		tx.Sender = request.Transaction.From
	}
	if request.Transaction.To != "" {
		tx.Recipient = request.Transaction.To
	}
	if request.Transaction.Nonce != "" {
		tx.Nonce = request.Transaction.Nonce
	}
	if request.Transaction.Value != "" {
		tx.Value = request.Transaction.Value
	}
	if request.Transaction.GasPrice != "" {
		tx.GasPrice = request.Transaction.GasPrice
	}
	if request.Transaction.GasLimit != "" {
		tx.GasLimit = request.Transaction.GasLimit
	}
	if request.Transaction.Data != "" {
		tx.Data = request.Transaction.Data
	}
	if request.Transaction.PrivateFrom != "" {
		tx.PrivateFrom = request.Transaction.PrivateFrom
	}
	if len(request.Transaction.PrivateFor) > 0 {
		tx.PrivateFor = request.Transaction.PrivateFor
	}
	if request.Transaction.PrivateFrom != "" {
		tx.PrivateFrom = request.Transaction.PrivateFrom
	}
	if request.Transaction.PrivacyGroupID != "" {
		tx.PrivacyGroupID = request.Transaction.PrivacyGroupID
	}
	if request.Transaction.Raw != "" {
		tx.Raw = request.Transaction.Raw
	}
}
