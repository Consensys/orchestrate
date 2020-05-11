package transactions

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/interfaces"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"
)

//go:generate mockgen -source=send_tx.go -destination=mocks/send_tx.go -package=mocks

const sendTxComponent = "use-cases.send-tx"

type SendTxUseCase interface {
	Execute(ctx context.Context, txRequest *types.TransactionRequest, tenantID string) (*types.TransactionResponse, error)
}

// sendTxUsecase is a use case to create a new transaction request
type sendTxUsecase struct {
	validator       validators.TransactionValidator
	db              interfaces.DB
	startJobUsecase jobs.StartJobUseCase
}

// NewSendTxUseCase creates a new SendTxUseCase
func NewSendTxUseCase(validator validators.TransactionValidator, db interfaces.DB, startJobUsecase jobs.StartJobUseCase) SendTxUseCase {
	return &sendTxUsecase{
		validator:       validator,
		db:              db,
		startJobUsecase: startJobUsecase,
	}
}

// Execute validates, creates and starts a new contract transaction
func (uc *sendTxUsecase) Execute(ctx context.Context, txRequest *types.TransactionRequest, tenantID string) (*types.TransactionResponse, error) {
	log.WithContext(ctx).WithField("idempotency_key", txRequest.IdempotencyKey).Debug("creating new transaction")

	// TODO: Validation of args given methodSignature

	requestHash, err := uc.validator.ValidateRequestHash(ctx, txRequest.Params, txRequest.IdempotencyKey)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}
	err = uc.validator.ValidateChainExists(ctx, txRequest.ChainUUID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// TODO: Craft "Data" field here with MethodSignature and Args
	txData := ""
	jsonParams, err := utils.ObjectToJSON(txRequest.Params)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Create model and insert in DB
	dbtx, err := uc.db.Begin()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	schedule := &models.Schedule{
		TenantID:  tenantID,
		ChainUUID: txRequest.ChainUUID,
	}
	err = dbtx.Schedule().Insert(ctx, schedule)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	job := &models.Job{
		ScheduleID: schedule.ID,
		Type:       types.JobConstantinopleTransaction,
		Transaction: &models.Transaction{
			Sender:    txRequest.Params.From,
			Recipient: txRequest.Params.To,
			Value:     txRequest.Params.Value,
			GasPrice:  txRequest.Params.GasPrice,
			GasLimit:  txRequest.Params.Gas,
			Data:      txData,
		},
		Labels: txRequest.Labels,
	}
	err = dbtx.Job().Insert(ctx, job)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}
	schedule.Jobs = []*models.Job{job}

	logModel := &models.Log{
		JobID:   job.ID,
		Status:  types.JobStatusCreated,
		Message: "Job created for contract transaction request",
	}
	err = dbtx.Log().Insert(ctx, logModel)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}
	job.Logs = []*models.Log{logModel}

	txRequestModel := &models.TransactionRequest{
		IdempotencyKey: txRequest.IdempotencyKey,
		ScheduleID:     schedule.ID,
		Schedule:       schedule,
		RequestHash:    requestHash,
		Params:         jsonParams,
	}
	err = dbtx.TransactionRequest().SelectOrInsert(ctx, txRequestModel)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	err = dbtx.Commit()
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Start first job of the schedule
	err = uc.startJobUsecase.Execute(ctx, txRequestModel.Schedule.Jobs[0].UUID, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	txResponse, err := utils.FormatTxResponse(txRequestModel)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	log.WithContext(ctx).WithField("idempotency_key", txRequestModel.IdempotencyKey).Info("contract transaction request created successfully")
	return txResponse, nil
}
