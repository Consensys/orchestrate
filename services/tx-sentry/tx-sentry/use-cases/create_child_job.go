package usecases

import (
	"context"
	"math/big"

	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx-scheduler"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

//go:generate mockgen -source=create_child_job.go -destination=mocks/create_child_job.go -package=mocks

const createChildJobComponent = "use-cases.create-child-job"

type CreateChildJobUseCase interface {
	Execute(ctx context.Context, job *entities.Job) (string, error)
}

// createChildJobUseCase is a use case to create a new transaction job
type createChildJobUseCase struct {
	txSchedulerClient txscheduler.TransactionSchedulerClient
}

// NewCreateChildJobUseCase creates a new StartSessionUseCase
func NewCreateChildJobUseCase(txSchedulerClient txscheduler.TransactionSchedulerClient) CreateChildJobUseCase {
	return &createChildJobUseCase{
		txSchedulerClient: txSchedulerClient,
	}
}

// Execute starts a job session
func (uc *createChildJobUseCase) Execute(ctx context.Context, job *entities.Job) (string, error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("verifying job status")

	jobs, err := uc.txSchedulerClient.SearchJob(ctx, &entities.JobFilters{
		ChainUUID:     job.ChainUUID,
		ParentJobUUID: job.UUID,
	})
	if err != nil {
		errMessage := "failed to get jobs"
		logger.Error(errMessage)
		return "", errors.FromError(err).ExtendComponent(createChildJobComponent)
	}

	parentJob := jobs[0]
	nJobs := float64(len(jobs))
	status := parentJob.Status
	if status != utils.StatusPending {
		logger.WithField("status", status).Info("job has been updated. stopping job session")
		return "", nil
	}

	gasPriceMultiplier := getGasPriceMultiplier(
		parentJob.Annotations.GasPricePolicy.RetryPolicy.Increment,
		parentJob.Annotations.GasPricePolicy.RetryPolicy.Limit,
		nJobs,
	)
	childJobRequest := newChildJobRequest(parentJob, gasPriceMultiplier)
	childJob, err := uc.txSchedulerClient.CreateJob(ctx, childJobRequest)
	if err != nil {
		errMessage := "failed create new child job"
		logger.Error(errMessage)
		return "", errors.FromError(err).ExtendComponent(createChildJobComponent)
	}

	err = uc.txSchedulerClient.StartJob(ctx, childJob.UUID)
	if err != nil {
		errMessage := "failed start child job"
		logger.WithField("child_job_uuid", childJob.UUID).Error(errMessage)
		return "", errors.FromError(err).ExtendComponent(createChildJobComponent)
	}

	logger.WithField("child_job_uuid", childJob.UUID).Info("new child job created")

	return childJob.UUID, nil
}

func getGasPriceMultiplier(increment, limit, nChildren float64) float64 {
	// This is fine as GasPriceIncrement default value is 0
	newGasPriceMultiplier := nChildren * increment

	if newGasPriceMultiplier >= limit {
		newGasPriceMultiplier = limit
	}

	return newGasPriceMultiplier
}

func newChildJobRequest(parentJob *txschedulertypes.JobResponse, gasPriceMultiplier float64) *txschedulertypes.CreateJobRequest {
	// We selectively choose fields from the parent job
	newJobRequest := &txschedulertypes.CreateJobRequest{
		ChainUUID:     parentJob.ChainUUID,
		ScheduleUUID:  parentJob.ScheduleUUID,
		Type:          parentJob.Type,
		Labels:        parentJob.Labels,
		Annotations:   parentJob.Annotations,
		ParentJobUUID: parentJob.UUID,
	}

	// raw transactions are resent as-is with no modifications
	if parentJob.Type == utils.EthereumRawTransaction {
		newJobRequest.Transaction = entities.ETHTransaction{
			Raw: parentJob.Transaction.Raw,
		}

		return newJobRequest
	}

	newJobRequest.Transaction = entities.ETHTransaction{
		From:           parentJob.Transaction.From,
		To:             parentJob.Transaction.To,
		Value:          parentJob.Transaction.Value,
		Data:           parentJob.Transaction.Data,
		PrivateFrom:    parentJob.Transaction.PrivateFrom,
		PrivateFor:     parentJob.Transaction.PrivateFor,
		PrivacyGroupID: parentJob.Transaction.PrivacyGroupID,
		Nonce:          parentJob.Transaction.Nonce,
	}

	gasPrice := new(big.Float)
	gasPrice, _ = gasPrice.SetString(parentJob.Transaction.GasPrice)
	newJobRequest.Transaction.GasPrice = gasPrice.Mul(gasPrice, big.NewFloat(1+gasPriceMultiplier)).String()

	return newJobRequest
}
