package transactions

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	tsorm "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/orm"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"
)

//go:generate mockgen -source=send_tx.go -destination=mocks/send_tx.go -package=mocks

const sendTxComponent = "use-cases.send-tx"

type SendTxUseCase interface {
	Execute(ctx context.Context, txRequest *types.TransactionRequest, chainUUID, tenantID string) (*types.TransactionResponse, error)
}

// sendTxUsecase is a use case to create a new transaction request
type sendTxUsecase struct {
	validator  validators.TransactionValidator
	db         store.DB
	startJobUC jobs.StartJobUseCase
	orm        tsorm.ORM
}

// NewSendTxUseCase creates a new SendTxUseCase
func NewSendTxUseCase(validator validators.TransactionValidator, db store.DB, orm tsorm.ORM, startJobUsecase jobs.StartJobUseCase) SendTxUseCase {
	return &sendTxUsecase{
		validator:  validator,
		db:         db,
		orm:        orm,
		startJobUC: startJobUsecase,
	}
}

// Execute validates, creates and starts a new contract transaction
func (uc *sendTxUsecase) Execute(ctx context.Context, txRequest *types.TransactionRequest, chainUUID, tenantID string) (*types.TransactionResponse, error) {
	logger := log.WithContext(ctx)

	logger.
		WithField("idempotency_key", txRequest.IdempotencyKey).
		Debug("creating new transaction")
	// TODO: Validation of args given methodSignature

	requestHash, err := uc.validator.ValidateRequestHash(ctx, chainUUID, txRequest.Params, txRequest.IdempotencyKey)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}
	err = uc.validator.ValidateChainExists(ctx, chainUUID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// TODO: Craft "Data" field here with MethodSignature and Args
	txData := ""
	jsonParams, err := utils.ObjectToJSON(txRequest.Params)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Step 1: Persist the job  an its linked entities
	job := &models.Job{
		Schedule: &models.Schedule{
			TenantID:  tenantID,
			ChainUUID: chainUUID,
		},
		Type: types.JobConstantinopleTransaction,
		Transaction: &models.Transaction{
			Sender:    txRequest.Params.From,
			Recipient: txRequest.Params.To,
			Value:     txRequest.Params.Value,
			GasPrice:  txRequest.Params.GasPrice,
			GasLimit:  txRequest.Params.Gas,
			Data:      txData,
		},
		Labels: txRequest.Labels,
		Logs: []*models.Log{
			{
				Status:  types.JobStatusCreated,
				Message: "Job created for contract transaction request",
			},
		},
	}

	txRequestModel := &models.TransactionRequest{
		IdempotencyKey: txRequest.IdempotencyKey,
		RequestHash:    requestHash,
		Params:         jsonParams,
	}

	err = database.ExecuteInDBTx(uc.db, func(tx database.Tx) error {
		if der := uc.orm.InsertOrUpdateJob(ctx, tx.(store.Tx), job); der != nil {
			return der
		}

		txRequestModel.ScheduleID = &job.Schedule.ID

		if der := tx.(store.Tx).TransactionRequest().SelectOrInsert(ctx, txRequestModel); der != nil {
			return der
		}

		return nil
	})

	if err != nil {
		logger.
			WithError(err).
			Errorf("cannot persist transaction request entity")
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	logger.
		WithField("job_uuid", job.UUID).
		Info("jbo was created successfully")

	// Step2: Start first job of the schedule
	err = uc.startJobUC.Execute(ctx, job.UUID, tenantID)
	if err != nil {
		logger.
			WithError(err).
			WithField("job_uuid", job.UUID).
			Errorf("cannot start job")
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Step3: Load latest Schedule status from DB
	schedule, err := uc.orm.FetchScheduleByID(ctx, uc.db, *txRequestModel.ScheduleID)
	if err != nil {
		logger.
			WithError(err).
			WithField("job_uuid", job.UUID).
			WithField("tenant_id", tenantID).
			Errorf("cannot find job")
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	txRequestModel.Schedule = schedule
	txResponse, err := utils.FormatTxResponse(txRequestModel)
	if err != nil {
		logger.
			WithError(err).
			Errorf("cannot format tx response")
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	logger.
		WithField("idempotency_key", txRequestModel.IdempotencyKey).
		Info("contract transaction request created successfully")
	return txResponse, nil
}
