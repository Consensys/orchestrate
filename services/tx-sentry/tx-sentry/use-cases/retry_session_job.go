package usecases

import (
	"context"
	"math"
	"math/big"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/log"
	orchestrateclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
)

//go:generate mockgen -source=retry_session_job.go -destination=mocks/retry_session_job.go -package=mocks

const retrySessionJobComponent = "use-cases.retry-session-job"

type RetrySessionJobUseCase interface {
	Execute(ctx context.Context, parentJobUUID, lastChildUUID string, nChildren int) (string, error)
}

// retrySessionJobUseCase is a use case to create a new transaction job
type retrySessionJobUseCase struct {
	client orchestrateclient.OrchestrateClient
	logger *log.Logger
}

// NewRetrySessionJobUseCase creates a new StartSessionUseCase
func NewRetrySessionJobUseCase(client orchestrateclient.OrchestrateClient) RetrySessionJobUseCase {
	return &retrySessionJobUseCase{
		client: client,
		logger: log.NewLogger().SetComponent(retrySessionJobComponent),
	}
}

// Execute starts a job session
func (uc *retrySessionJobUseCase) Execute(ctx context.Context, jobUUID, childUUID string, nChildren int) (string, error) {
	logger := uc.logger.WithContext(ctx).WithField("job", jobUUID)
	ctx = log.With(ctx, logger)

	logger.Debug("verifying job status")

	job, err := uc.client.GetJob(ctx, jobUUID)
	if err != nil {
		logger.WithError(err).Error("failed to get job")
		return "", errors.FromError(err).ExtendComponent(retrySessionJobComponent)
	}

	status := job.Status
	if status != utils.StatusPending {
		logger.WithField("status", status).Info("job has been updated. stopping job session")
		return "", nil
	}

	// In case gas increments on every retry we create a new job
	if job.Type != utils.EthereumRawTransaction &&
		(job.Annotations.GasPricePolicy.RetryPolicy.Increment > 0.0 &&
			nChildren <= int(math.Ceil(job.Annotations.GasPricePolicy.RetryPolicy.Limit/job.Annotations.GasPricePolicy.RetryPolicy.Increment))) {

		childJob, errr := uc.createAndStartNewChildJob(ctx, job, nChildren)
		if errr != nil {
			return "", errors.FromError(errr).ExtendComponent(retrySessionJobComponent)
		}

		logger.WithField("child_job", childJob.UUID).Info("new child job created and started")
		return childJob.UUID, nil
	}

	// Otherwise we retry on last job
	err = uc.client.ResendJobTx(ctx, childUUID)
	if err != nil {
		logger.WithError(err).Error("failed to resend job")
		return "", errors.FromError(err).ExtendComponent(retrySessionJobComponent)
	}

	logger.Info("job has been resent")
	return job.UUID, nil
}

func (uc *retrySessionJobUseCase) createAndStartNewChildJob(ctx context.Context,
	parentJob *types.JobResponse,
	nChildrenJobs int,
) (*types.JobResponse, error) {
	logger := uc.logger.WithContext(ctx).WithField("job", parentJob.UUID)
	gasPriceMultiplier := getGasPriceMultiplier(
		parentJob.Annotations.GasPricePolicy.RetryPolicy.Increment,
		parentJob.Annotations.GasPricePolicy.RetryPolicy.Limit,
		float64(nChildrenJobs),
	)

	childJobRequest := newChildJobRequest(parentJob, gasPriceMultiplier)
	childJob, err := uc.client.CreateJob(ctx, childJobRequest)
	if err != nil {
		logger.Error("failed create new child job")
		return nil, errors.FromError(err).ExtendComponent(retrySessionJobComponent)
	}

	err = uc.client.StartJob(ctx, childJob.UUID)
	if err != nil {
		logger.WithField("child_job", childJob.UUID).Error("failed start child job")
		return nil, errors.FromError(err).ExtendComponent(retrySessionJobComponent)
	}

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

func newChildJobRequest(parentJob *types.JobResponse, gasPriceMultiplier float64) *types.CreateJobRequest {
	// We selectively choose fields from the parent job
	newJobRequest := &types.CreateJobRequest{
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
