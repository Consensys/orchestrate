package transactions

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"
)

//go:generate mockgen -source=send_tx.go -destination=mocks/send_tx.go -package=mocks

const sendTxComponent = "use-cases.send-tx"

type SendTxUseCase interface {
	Execute(ctx context.Context, txRequest *entities.TxRequest, chainUUID, tenantID string) (*entities.TxRequest, error)
}

// sendTxUsecase is a use case to create a new transaction request
type sendTxUsecase struct {
	validator        validators.TransactionValidator
	db               store.DB
	startJobUC       jobs.StartJobUseCase
	createJobUC      jobs.CreateJobUseCase
	createScheduleUC schedules.CreateScheduleUseCase
	getScheduleUC    schedules.GetScheduleUseCase
}

// NewSendTxUseCase creates a new SendTxUseCase
func NewSendTxUseCase(validator validators.TransactionValidator,
	db store.DB,
	startJobUseCase jobs.StartJobUseCase,
	createJobUC jobs.CreateJobUseCase,
	createScheduleUC schedules.CreateScheduleUseCase,
	getScheduleUC schedules.GetScheduleUseCase,
) SendTxUseCase {
	return &sendTxUsecase{
		validator:        validator,
		db:               db,
		startJobUC:       startJobUseCase,
		createJobUC:      createJobUC,
		createScheduleUC: createScheduleUC,
		getScheduleUC:    getScheduleUC,
	}
}

// Execute validates, creates and starts a new contract transaction
func (uc *sendTxUsecase) Execute(ctx context.Context, txRequest *entities.TxRequest, chainUUID, tenantID string) (*entities.TxRequest, error) {
	logger := log.WithContext(ctx)

	logger.
		WithField("idempotency_key", txRequest.IdempotencyKey).
		Debug("creating new transaction")

	if err := utils.GetValidator().Struct(txRequest); err != nil {
		return nil, errors.InvalidParameterError(err.Error()).ExtendComponent(sendTxComponent)
	}

	if err := txRequest.Params.PrivateTransactionParams.Validate(); err != nil {
		return nil, errors.InvalidParameterError(err.Error()).ExtendComponent(sendTxComponent)
	}

	// Step 1: Validate RequestHash
	requestHash, err := uc.validator.ValidateRequestHash(ctx, chainUUID, txRequest.Params, txRequest.IdempotencyKey)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	txDataBytes, err := uc.validator.ValidateMethodSignature(txRequest.Params.MethodSignature, txRequest.Params.Args)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	jobType := generateJobType(txRequest)

	// Step 2: Insert Schedule + Job + Transaction + TxRequest atomically
	err = database.ExecuteInDBTx(uc.db, func(dbtx database.Tx) error {
		txRequestModel, der := parsers.NewTxRequestModelFromEntities(txRequest, requestHash)
		if der != nil {
			return der
		}

		der = dbtx.(store.Tx).TransactionRequest().SelectOrInsert(ctx, txRequestModel)
		if der != nil {
			return der
		}
		txRequest.CreatedAt = txRequestModel.CreatedAt

		txRequest.Schedule, der = uc.createScheduleUC.
			WithDBTransaction(dbtx.(store.Tx)).
			Execute(ctx, &entities.Schedule{}, tenantID)

		if der != nil {
			return der
		}

		// Craft "Data" field
		sendTxJob := parsers.NewJobEntityFromTxRequest(txRequest, jobType, chainUUID)
		sendTxJob.Transaction.Data = hexutil.Encode(txDataBytes)

		job, der := uc.createJobUC.WithDBTransaction(dbtx.(store.Tx)).Execute(ctx, sendTxJob, tenantID)
		if der != nil {
			return der
		}

		txRequest.Schedule.Jobs = []*types.Job{job}

		return nil
	})

	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Step 3: Start first job of the schedule
	scheduleUUID := txRequest.Schedule.UUID
	jobUUID := txRequest.Schedule.Jobs[0].UUID
	err = uc.startJobUC.Execute(ctx, jobUUID, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Step 4: Load latest Schedule status from DB
	txRequest.Schedule, err = uc.getScheduleUC.Execute(ctx, scheduleUUID, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	logger.
		WithField("idempotency_key", txRequest.IdempotencyKey).
		WithField("schedule_uuid", scheduleUUID).
		WithField("job_uuid", jobUUID).
		Info("contract transaction request created successfully")

	return txRequest, nil
}

func generateJobType(txRequest *entities.TxRequest) string {
	switch {
	case txRequest.Params.Protocol == utils.OrionChainType:
		return types.OrionEEATransaction
	case txRequest.Params.Protocol == utils.TesseraChainType:
		return types.TesseraPrivateTransaction
	default:
		return types.EthereumTransaction
	}
}
