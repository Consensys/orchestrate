package transactions

import (
	"context"

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
	Execute(ctx context.Context, txRequest *entities.TxRequest, txData, chainUUID, tenantID string) (*entities.TxRequest, error)
}

// sendTxUsecase is a use case to create a new transaction
type sendTxUsecase struct {
	validator        validators.TransactionValidator
	db               store.DB
	startJobUC       jobs.StartJobUseCase
	createJobUC      jobs.CreateJobUseCase
	createScheduleUC schedules.CreateScheduleUseCase
	getTxUC          GetTxUseCase
}

// NewSendTxUseCase creates a new SendTxUseCase
func NewSendTxUseCase(validator validators.TransactionValidator,
	db store.DB,
	startJobUseCase jobs.StartJobUseCase,
	createJobUC jobs.CreateJobUseCase,
	createScheduleUC schedules.CreateScheduleUseCase,
	getTxUC GetTxUseCase,
) SendTxUseCase {
	return &sendTxUsecase{
		validator:        validator,
		db:               db,
		startJobUC:       startJobUseCase,
		createJobUC:      createJobUC,
		createScheduleUC: createScheduleUC,
		getTxUC:          getTxUC,
	}
}

// Execute validates, creates and starts a new transaction
func (uc *sendTxUsecase) Execute(ctx context.Context, txRequest *entities.TxRequest, txData, chainUUID, tenantID string) (*entities.TxRequest, error) {
	logger := log.WithContext(ctx).WithField("idempotency_key", txRequest.IdempotencyKey)
	logger.Debug("creating new transaction")

	// Step 1: Validation
	err := uc.validator.ValidateFields(ctx, txRequest)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	requestHash, err := uc.validator.ValidateRequestHash(ctx, chainUUID, txRequest.Params, txRequest.IdempotencyKey)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Step 2: Insert Schedule + Job + Transaction + TxRequest atomically
	err = database.ExecuteInDBTx(uc.db, func(dbtx database.Tx) error {
		txRequestModel := parsers.NewTxRequestModelFromEntities(txRequest, requestHash)

		der := dbtx.(store.Tx).TransactionRequest().SelectOrInsert(ctx, txRequestModel)
		if der != nil {
			return der
		}
		txRequest.UUID = txRequestModel.UUID

		txRequest.Schedule, der = uc.createScheduleUC.
			WithDBTransaction(dbtx.(store.Tx)).
			Execute(ctx, &entities.Schedule{TxRequest: txRequest}, tenantID)

		if der != nil {
			return der
		}

		sendTxJob := parsers.NewJobEntityFromTxRequest(txRequest, generateJobType(txRequest), chainUUID)
		sendTxJob.Transaction.Data = txData

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
	jobUUID := txRequest.Schedule.Jobs[0].UUID
	err = uc.startJobUC.Execute(ctx, jobUUID, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Step 4: Load latest Schedule status from DB
	txRequest, err = uc.getTxUC.Execute(ctx, txRequest.UUID, []string{tenantID})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	logger.WithField("uuid", txRequest.UUID).Info("contract transaction request created successfully")
	return txRequest, nil
}

func generateJobType(txRequest *entities.TxRequest) string {
	switch {
	case txRequest.Params.Protocol == utils.OrionChainType:
		return types.OrionEEATransaction
	case txRequest.Params.Protocol == utils.TesseraChainType:
		return types.TesseraPrivateTransaction
	case txRequest.Params.Raw != "":
		return types.EthereumRawTransaction
	default:
		return types.EthereumTransaction
	}
}
