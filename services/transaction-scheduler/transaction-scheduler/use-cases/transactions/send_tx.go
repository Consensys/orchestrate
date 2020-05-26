package transactions

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
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
	Execute(ctx context.Context, txRequest *entities.TxRequest, tenantID string) (*entities.TxRequest, error)
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
func (uc *sendTxUsecase) Execute(ctx context.Context, txRequest *entities.TxRequest, tenantID string) (*entities.TxRequest, error) {
	logger := log.WithContext(ctx)

	logger.
		WithField("idempotency_key", txRequest.IdempotencyKey).
		Debug("creating new transaction")

	// Step 1: Validate RequestHash and ChainUUID
	// @TODO Validation of args given methodSignature
	requestHash, err := uc.validator.ValidateRequestHash(ctx, txRequest.Schedule.ChainUUID, txRequest.Params,
		txRequest.IdempotencyKey)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	txRequestModel, err := parsers.NewTxRequestModelFromEntities(txRequest, requestHash, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Step 2: Insert Schedule + Job + Transaction + TxRequest atomically
	err = database.ExecuteInDBTx(uc.db, func(dbtx database.Tx) error {
		der := dbtx.(store.Tx).TransactionRequest().SelectOrInsert(ctx, txRequestModel)
		if der != nil {
			return der
		}

		txRequest.Schedule, der = uc.createScheduleUC.
			WithDBTransaction(dbtx.(store.Tx)).
			Execute(ctx, txRequest.Schedule, tenantID)
		if der != nil {
			return der
		}

		sendTxJob, der := parsers.NewJobEntityFromSendTxRequest(txRequest)
		if der != nil {
			return der
		}

		txRequest.Schedule.Jobs[0], der = uc.createJobUC.
			WithDBTransaction(dbtx.(store.Tx)).
			Execute(ctx, sendTxJob, tenantID)
		if der != nil {
			return der
		}

		return nil
	})

	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	scheduleUUID := txRequest.Schedule.UUID

	// Step 3: Start first job of the schedule
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
