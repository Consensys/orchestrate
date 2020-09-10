package usecases

import (
	"context"
	"fmt"
	"math"
	"math/big"

	txschedulertypes "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/txscheduler"
	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

//go:generate mockgen -source=retry_session_job.go -destination=mocks/retry_session_job.go -package=mocks

const retrySessionJobComponent = "use-cases.retry-session-job"

type RetrySessionJobUseCase interface {
	Execute(ctx context.Context, job *entities.Job) (string, error)
}

// retrySessionJobUseCase is a use case to create a new transaction job
type retrySessionJobUseCase struct {
	txSchedulerClient txscheduler.TransactionSchedulerClient
}

// NewRetrySessionJobUseCase creates a new StartSessionUseCase
func NewRetrySessionJobUseCase(txSchedulerClient txscheduler.TransactionSchedulerClient) RetrySessionJobUseCase {
	return &retrySessionJobUseCase{
		txSchedulerClient: txSchedulerClient,
	}
}

// Execute starts a job session
func (uc *retrySessionJobUseCase) Execute(ctx context.Context, job *entities.Job) (string, error) {
	logger := log.WithContext(ctx).WithField("job_uuid", job.UUID)
	logger.Debug("verifying job status")

	jobs, err := uc.txSchedulerClient.SearchJob(ctx, &entities.JobFilters{
		ChainUUID:     job.ChainUUID,
		ParentJobUUID: job.UUID,
	})
	if err != nil {
		errMessage := "failed to get jobs"
		logger.Error(errMessage)
		return "", errors.FromError(err).ExtendComponent(retrySessionJobComponent)
	}

	parentJob := jobs[0]
	lastJobRetry := jobs[len(jobs)-1]
	status := parentJob.Status
	if status != utils.StatusPending {
		logger.WithField("status", status).Info("job has been updated. stopping job session")
		return "", nil
	}

	// we count the number of resending of last job as retries
	nRetries := len(jobs) - 1
	for _, lg := range lastJobRetry.Logs {
		if lg.Status == utils.StatusResending {
			nRetries++
		}
	}

	if nRetries >= txschedulertypes.SentryMaxRetries {
		logger.WithField("retries", nRetries).Infof("job session exceeded max amount of retries %d", nRetries)
		return "", nil
	}

	childrenJobs := len(jobs) - 1
	// In case gas increments on every retry we create a new job
	if parentJob.Annotations.GasPricePolicy.RetryPolicy.Increment > 0.0 &&
		childrenJobs <= int(math.Ceil(parentJob.Annotations.GasPricePolicy.RetryPolicy.Limit/parentJob.Annotations.GasPricePolicy.RetryPolicy.Increment)) {

		childJob, errr := uc.CreateAndStartNewChildJob(ctx, parentJob, childrenJobs)
		if errr != nil {
			return "", errors.FromError(errr).ExtendComponent(retrySessionJobComponent)
		}

		return childJob.UUID, nil
	}

	// Otherwise we retry on last job
	err = uc.txSchedulerClient.ResendJobTx(ctx, lastJobRetry.UUID)
	if err != nil {
		return "", errors.FromError(err).ExtendComponent(retrySessionJobComponent)
	}

	return parentJob.UUID, nil
}

func (uc *retrySessionJobUseCase) CreateAndStartNewChildJob(ctx context.Context,
	parentJob *txschedulertypes.JobResponse,
	nChildrenJobs int,
) (*txschedulertypes.JobResponse, error) {
	logger := log.WithContext(ctx).WithField("job_uuid", parentJob.UUID)
	gasPriceMultiplier := getGasPriceMultiplier(
		parentJob.Annotations.GasPricePolicy.RetryPolicy.Increment,
		parentJob.Annotations.GasPricePolicy.RetryPolicy.Limit,
		float64(nChildrenJobs),
	)

	childJobRequest := newChildJobRequest(parentJob, gasPriceMultiplier, nChildrenJobs)
	childJob, err := uc.txSchedulerClient.CreateJob(ctx, childJobRequest)
	if err != nil {
		errMessage := "failed create new child job"
		logger.Error(errMessage)
		return nil, errors.FromError(err).ExtendComponent(retrySessionJobComponent)
	}

	err = uc.txSchedulerClient.StartJob(ctx, childJob.UUID)
	if err != nil {
		errMessage := "failed start child job"
		logger.WithField("child_job_uuid", childJob.UUID).Error(errMessage)
		return nil, errors.FromError(err).ExtendComponent(retrySessionJobComponent)
	}

	logger.WithField("child_job_uuid", childJob.UUID).Info("new child job created")

	return childJob, nil
}

func getGasPriceMultiplier(increment, limit, nChildren float64) float64 {
	// This is fine as GasPriceIncrement default value is 0
	newGasPriceMultiplier := (nChildren + 1) * increment

	if newGasPriceMultiplier >= limit {
		newGasPriceMultiplier = limit
	}

	return newGasPriceMultiplier
}

func newChildJobRequest(parentJob *txschedulertypes.JobResponse, gasPriceMultiplier float64, nRetries int) *txschedulertypes.CreateJobRequest {
	// We selectively choose fields from the parent job
	newJobRequest := &txschedulertypes.CreateJobRequest{
		ChainUUID:     parentJob.ChainUUID,
		ScheduleUUID:  parentJob.ScheduleUUID,
		Type:          parentJob.Type,
		Labels:        parentJob.Labels,
		Annotations:   parentJob.Annotations,
		ParentJobUUID: parentJob.UUID,
	}

	newJobRequest.Labels["retryOrder"] = fmt.Sprintf("%d", nRetries+1)

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
