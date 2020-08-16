package transactions

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"

	chainregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/parsers"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"
)

//go:generate mockgen -source=send_tx.go -destination=mocks/send_tx.go -package=mocks

const sendTxComponent = "use-cases.send-tx"

type SendTxUseCase interface {
	Execute(ctx context.Context, txRequest *entities.TxRequest, txData, tenantID string) (*entities.TxRequest, error)
}

// sendTxUsecase is a use case to create a new transaction
type sendTxUsecase struct {
	validator           validators.TransactionValidator
	db                  store.DB
	chainRegistryCLient chainregistry.ChainRegistryClient
	startJobUC          jobs.StartJobUseCase
	createJobUC         jobs.CreateJobUseCase
	createScheduleUC    schedules.CreateScheduleUseCase
	getTxUC             GetTxUseCase
}

// NewSendTxUseCase creates a new SendTxUseCase
func NewSendTxUseCase(
	validator validators.TransactionValidator,
	db store.DB,
	chainRegistryCLient chainregistry.ChainRegistryClient,
	startJobUseCase jobs.StartJobUseCase,
	createJobUC jobs.CreateJobUseCase,
	createScheduleUC schedules.CreateScheduleUseCase,
	getTxUC GetTxUseCase,
) SendTxUseCase {
	return &sendTxUsecase{
		validator:           validator,
		db:                  db,
		chainRegistryCLient: chainRegistryCLient,
		startJobUC:          startJobUseCase,
		createJobUC:         createJobUC,
		createScheduleUC:    createScheduleUC,
		getTxUC:             getTxUC,
	}
}

// Execute validates, creates and starts a new transaction
func (uc *sendTxUsecase) Execute(ctx context.Context, txRequest *entities.TxRequest, txData, tenantID string) (*entities.TxRequest, error) {
	logger := log.WithContext(ctx).WithField("idempotency_key", txRequest.IdempotencyKey)
	logger.Debug("creating new transaction")

	// Step 1: Get chainUUID from chain registry
	chainUUID, err := uc.getChainUUID(ctx, txRequest.ChainName)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Step 2: Generate request hash
	requestHash, err := generateRequestHash(chainUUID, txRequest.Params)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Step 3: Insert Schedule + Job + Transaction + TxRequest atomically OR get tx request if it exists
	txRequest, err = uc.selectOrInsertTxRequest(ctx, txRequest, txData, requestHash, chainUUID, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Step 4: Start first job of the schedule if status is CREATED
	job := txRequest.Schedule.Jobs[0]
	if job.GetStatus() == utils.StatusCreated {
		err = uc.startJobUC.Execute(ctx, job.UUID, []string{tenantID})
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
		}
	}

	// Step 5: Load latest Schedule status from DB
	txRequest, err = uc.getTxUC.Execute(ctx, txRequest.UUID, []string{tenantID})
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	logger.WithField("uuid", txRequest.UUID).Info("send transaction request created successfully")
	return txRequest, nil
}

func (uc *sendTxUsecase) getChainUUID(ctx context.Context, chainName string) (string, error) {
	chain, err := uc.chainRegistryCLient.GetChainByName(ctx, chainName)
	if err != nil {
		errMessage := fmt.Sprintf("cannot load '%s' chain", chainName)
		log.WithContext(ctx).WithError(err).Error(errMessage)
		if errors.IsNotFoundError(err) {
			return "", errors.InvalidParameterError(errMessage)
		}
		return "", errors.FromError(err)
	}

	return chain.UUID, nil
}

func (uc *sendTxUsecase) selectOrInsertTxRequest(
	ctx context.Context,
	txRequest *entities.TxRequest,
	txData, requestHash, chainUUID, tenantID string,
) (*entities.TxRequest, error) {
	txRequestModel, err := uc.db.TransactionRequest().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID)
	switch {
	case errors.IsNotFoundError(err):
		return uc.insertNewTxRequest(ctx, txRequest, txData, requestHash, chainUUID, tenantID)
	case err != nil:
		return nil, err
	case txRequestModel != nil && txRequestModel.RequestHash != requestHash:
		errMessage := "a transaction request with the same idempotency key and different params already exists"
		log.WithError(err).WithField("idempotency_key", txRequestModel.IdempotencyKey).Error(errMessage)
		return nil, errors.AlreadyExistsError(errMessage)
	default:
		return uc.getTxUC.Execute(ctx, txRequestModel.UUID, []string{tenantID})
	}
}

func (uc *sendTxUsecase) insertNewTxRequest(
	ctx context.Context,
	txRequest *entities.TxRequest,
	txData, requestHash, chainUUID, tenantID string,
) (*entities.TxRequest, error) {
	err := database.ExecuteInDBTx(uc.db, func(dbtx database.Tx) error {
		schedule, der := uc.createScheduleUC.WithDBTransaction(dbtx.(store.Tx)).Execute(ctx, &entities.Schedule{TenantID: tenantID})
		if der != nil {
			return der
		}
		txRequest.Schedule = schedule

		scheduleModel, der := dbtx.(store.Tx).Schedule().FindOneByUUID(ctx, txRequest.Schedule.UUID, []string{tenantID})
		if der != nil {
			return der
		}

		txRequestModel := parsers.NewTxRequestModelFromEntities(txRequest, requestHash, scheduleModel.ID)
		der = dbtx.(store.Tx).TransactionRequest().Insert(ctx, txRequestModel)
		if der != nil {
			return der
		}
		txRequest.UUID = txRequestModel.UUID

		sendTxJobs := parsers.NewJobEntitiesFromTxRequest(txRequest, chainUUID, txData)
		txRequest.Schedule.Jobs = make([]*entities.Job, len(sendTxJobs))
		var nextJobUUID string
		for idx, txJob := range sendTxJobs {
			if nextJobUUID != "" {
				txJob.UUID = nextJobUUID
			}

			if idx < len(sendTxJobs)-1 {
				nextJobUUID = uuid.Must(uuid.NewV4()).String()
				txJob.NextJobUUID = nextJobUUID
			}

			job, der := uc.createJobUC.WithDBTransaction(dbtx.(store.Tx)).Execute(ctx, txJob, []string{tenantID})
			if der != nil {
				return der
			}

			txRequest.Schedule.Jobs[idx] = job
		}

		return nil
	})

	return txRequest, err
}

func generateRequestHash(chainUUID string, params interface{}) (string, error) {
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	hash := md5.Sum([]byte(string(jsonParams) + chainUUID))
	return hex.EncodeToString(hash[:]), nil
}
