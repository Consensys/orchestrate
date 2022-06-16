package jobs

import (
	"context"
	"math"
	"math/big"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/services/api/business/parsers"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
)

const retryJobTxComponent = "use-cases.retry-job-tx"

type retryJobTxUseCase struct {
	db            store.DB
	createJobTxUC usecases.CreateJobUseCase
	startJobUC    usecases.StartJobUseCase
	logger        *log.Logger
}

func NewRetryJobTxUseCase(db store.DB, createJobTxUC usecases.CreateJobUseCase, startJobUC usecases.StartJobUseCase) usecases.RetryJobTxUseCase {
	return &retryJobTxUseCase{
		db:            db,
		createJobTxUC: createJobTxUC,
		startJobUC:    startJobUC,
		logger:        log.NewLogger().SetComponent(retryJobTxComponent),
	}
}

// Execute sends a job to the Kafka topic
func (uc *retryJobTxUseCase) Execute(ctx context.Context, jobUUID string, gasIncrement float64, txData hexutil.Bytes, userInfo *multitenancy.UserInfo) error {
	ctx = log.WithFields(ctx, log.Field("job", jobUUID))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("retrying job transaction")

	jobModel, err := uc.db.Job().FindOneByUUID(ctx, jobUUID, userInfo.AllowedTenants, userInfo.Username, false)
	if err != nil {
		return errors.FromError(err).ExtendComponent(retryJobTxComponent)
	}

	job := parsers.NewJobEntityFromModels(jobModel)
	if job.Status != entities.StatusPending {
		errMessage := "cannot retry job transaction at the current status"
		logger.WithField("status", job.Status).Error(errMessage)
		return errors.InvalidStateError(errMessage).ExtendComponent(retryJobTxComponent)
	}

	job.UUID = ""
	job.InternalData.ParentJobUUID = jobUUID
	job.Transaction.Data = txData
	increment := int64(math.Trunc((gasIncrement + 1.0) * 100))
	if job.Transaction.TransactionType == entities.LegacyTxType {
		gasPrice := job.Transaction.GasPrice.ToInt()
		txGasPrice := gasPrice.Mul(gasPrice, big.NewInt(increment)).Div(gasPrice, big.NewInt(100))
		job.Transaction.GasPrice = utils.ToPtr(hexutil.Big(*txGasPrice)).(*hexutil.Big)
	} else {
		gasFeeCap := job.Transaction.GasFeeCap.ToInt()
		txGasFeeCap := gasFeeCap.Mul(gasFeeCap, big.NewInt(increment)).Div(gasFeeCap, big.NewInt(100))
		job.Transaction.GasFeeCap = utils.ToPtr(hexutil.Big(*txGasFeeCap)).(*hexutil.Big)
	}

	retriedJob, err := uc.createJobTxUC.Execute(ctx, job, userInfo)
	if err != nil {
		return errors.FromError(err).ExtendComponent(retryJobTxComponent)
	}

	if err = uc.startJobUC.Execute(ctx, retriedJob.UUID, userInfo); err != nil {
		return errors.FromError(err).ExtendComponent(retryJobTxComponent)
	}

	logger.WithField("job", retriedJob.UUID).Info("job retried successfully")
	return nil
}
