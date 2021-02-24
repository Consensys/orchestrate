package transactions

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/database"
	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/log"
	"github.com/ConsenSys/orchestrate/pkg/multitenancy"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/services/api/business/parsers"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/ConsenSys/orchestrate/services/api/store"
	"github.com/ConsenSys/orchestrate/services/api/store/models"
	"github.com/gofrs/uuid"
)

const sendTxComponent = "use-cases.send-tx"

// sendTxUsecase is a use case to create a new transaction
type sendTxUsecase struct {
	db                 store.DB
	searchChainsUC     usecases.SearchChainsUseCase
	startJobUC         usecases.StartJobUseCase
	createJobUC        usecases.CreateJobUseCase
	getTxUC            usecases.GetTxUseCase
	getFaucetCandidate usecases.GetFaucetCandidateUseCase
	logger             *log.Logger
}

// NewSendTxUseCase creates a new SendTxUseCase
func NewSendTxUseCase(
	db store.DB,
	searchChainsUC usecases.SearchChainsUseCase,
	startJobUseCase usecases.StartJobUseCase,
	createJobUC usecases.CreateJobUseCase,
	getTxUC usecases.GetTxUseCase,
	getFaucetCandidate usecases.GetFaucetCandidateUseCase,
) usecases.SendTxUseCase {
	return &sendTxUsecase{
		db:                 db,
		searchChainsUC:     searchChainsUC,
		startJobUC:         startJobUseCase,
		createJobUC:        createJobUC,
		getTxUC:            getTxUC,
		getFaucetCandidate: getFaucetCandidate,
		logger:             log.NewLogger().SetComponent(sendTxComponent),
	}
}

// Execute validates, creates and starts a new transaction
func (uc *sendTxUsecase) Execute(ctx context.Context, txRequest *entities.TxRequest, txData, tenantID string) (*entities.TxRequest, error) {
	ctx = log.WithFields(ctx, log.Field("idempotency-key", txRequest.IdempotencyKey))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("creating new transaction")

	allowedTenants := []string{tenantID, multitenancy.DefaultTenant}

	// Step 1: Get chain from chain registry
	chain, err := uc.getChain(ctx, txRequest.ChainName, allowedTenants)
	if err != nil {
		logger.WithError(err).WithField("chain_name", txRequest.ChainName).Error("failed to get chain")
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Step 2: Generate request hash
	requestHash, err := generateRequestHash(chain.UUID, txRequest.Params)
	if err != nil {
		logger.WithError(err).Error("failed to generate request hash")
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Step 3: Insert Schedule + Job + Transaction + TxRequest atomically OR get tx request if it exists
	txRequest, err = uc.selectOrInsertTxRequest(ctx, txRequest, txData, requestHash, chain.UUID, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Step 4: Start first job of the schedule if status is CREATED
	// Otherwise there was another request with same idempotency key and same reqHash
	job := txRequest.Schedule.Jobs[0]
	if job.Status == entities.StatusCreated {
		var fctJob *entities.Job
		fctJob, err = uc.startFaucetJob(ctx, txRequest.Params.From, job.ScheduleUUID, tenantID, chain)
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
		}
		if fctJob != nil {
			txRequest.Schedule.Jobs = append(txRequest.Schedule.Jobs, fctJob)
		}

		if err = uc.startJobUC.Execute(ctx, job.UUID, allowedTenants); err != nil {
			return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
		}
	} else { // Load latest Schedule status from DB
		txRequest, err = uc.getTxUC.Execute(ctx, txRequest.Schedule.UUID, []string{tenantID})
		if err != nil {
			return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
		}
	}

	logger.WithField("schedule", txRequest.Schedule.UUID).Info("transaction created successfully")
	return txRequest, nil
}

func (uc *sendTxUsecase) getChain(ctx context.Context, chainName string, tenants []string) (*entities.Chain, error) {
	chains, err := uc.searchChainsUC.Execute(ctx, &entities.ChainFilters{Names: []string{chainName}}, tenants)
	if err != nil {
		return nil, errors.FromError(err)
	}

	if len(chains) == 0 {
		errMessage := fmt.Sprintf("chain '%s' does not exist", chainName)
		return nil, errors.InvalidParameterError(errMessage)
	}

	return chains[0], nil
}

func (uc *sendTxUsecase) selectOrInsertTxRequest(
	ctx context.Context,
	txRequest *entities.TxRequest,
	txData, requestHash, chainUUID, tenantID string,
) (*entities.TxRequest, error) {
	if txRequest.IdempotencyKey == "" {
		return uc.insertNewTxRequest(ctx, txRequest, txData, requestHash, chainUUID, tenantID)
	}

	txRequestModel, err := uc.db.TransactionRequest().FindOneByIdempotencyKey(ctx, txRequest.IdempotencyKey, tenantID)
	switch {
	case errors.IsNotFoundError(err):
		return uc.insertNewTxRequest(ctx, txRequest, txData, requestHash, chainUUID, tenantID)
	case err != nil:
		return nil, err
	case txRequestModel != nil && txRequestModel.RequestHash != requestHash:
		errMessage := "transaction request with the same idempotency key and different params already exists"
		uc.logger.Error(errMessage)
		return nil, errors.AlreadyExistsError(errMessage)
	default:
		return uc.getTxUC.Execute(ctx, txRequestModel.Schedule.UUID, []string{tenantID})
	}
}

func (uc *sendTxUsecase) insertNewTxRequest(
	ctx context.Context,
	txRequest *entities.TxRequest,
	txData, requestHash, chainUUID, tenantID string,
) (*entities.TxRequest, error) {
	err := database.ExecuteInDBTx(uc.db, func(dbtx database.Tx) error {
		schedule := &models.Schedule{TenantID: tenantID}
		if err := dbtx.(store.Tx).Schedule().Insert(ctx, schedule); err != nil {
			return err
		}

		txRequestModel := parsers.NewTxRequestModelFromEntities(txRequest, requestHash, schedule.ID)
		if err := dbtx.(store.Tx).TransactionRequest().Insert(ctx, txRequestModel); err != nil {
			return err
		}

		txRequest.Schedule = parsers.NewScheduleEntityFromModels(schedule)
		return nil
	})

	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	sendTxJobs := parsers.NewJobEntitiesFromTxRequest(txRequest, chainUUID, txData)
	txRequest.Schedule.Jobs = make([]*entities.Job, len(sendTxJobs))
	err = database.ExecuteInDBTx(uc.db, func(dbtx database.Tx) error {
		var nextJobUUID string
		for idx, txJob := range sendTxJobs {
			if nextJobUUID != "" {
				txJob.UUID = nextJobUUID
			}

			if idx < len(sendTxJobs)-1 {
				nextJobUUID = uuid.Must(uuid.NewV4()).String()
				txJob.NextJobUUID = nextJobUUID
			}

			var job *entities.Job
			job, err = uc.createJobUC.WithDBTransaction(dbtx.(store.Tx)).
				Execute(ctx, txJob, []string{tenantID, multitenancy.DefaultTenant})
			if err != nil {
				return err
			}

			txRequest.Schedule.Jobs[idx] = job
		}
		return nil
	})

	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	return txRequest, nil
}

// Execute validates, creates and starts a new transaction for pre funding users account
func (uc *sendTxUsecase) startFaucetJob(ctx context.Context, account, scheduleUUID, tenantID string, chain *entities.Chain) (*entities.Job, error) {
	if account == "" {
		return nil, nil
	}

	logger := uc.logger.WithContext(ctx).WithField("chain", chain.UUID)
	faucet, err := uc.getFaucetCandidate.Execute(ctx, account, chain, []string{tenantID, multitenancy.DefaultTenant})
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, nil
		}
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}
	logger.WithField("faucet_amount", faucet.Amount).Debug("faucet: credit approved")

	txJob := &entities.Job{
		ScheduleUUID: scheduleUUID,
		ChainUUID:    chain.UUID,
		Type:         entities.EthereumTransaction,
		Labels: map[string]string{
			"faucetUUID": faucet.UUID,
		},
		InternalData: &entities.InternalData{},
		Transaction: &entities.ETHTransaction{
			From:  faucet.CreditorAccount,
			To:    account,
			Value: faucet.Amount,
		},
	}
	fctJob, err := uc.createJobUC.Execute(ctx, txJob, []string{tenantID, multitenancy.DefaultTenant})
	if err != nil {
		return nil, err
	}

	err = uc.startJobUC.Execute(ctx, fctJob.UUID, []string{tenantID, multitenancy.DefaultTenant})
	if err != nil {
		return fctJob, err
	}

	return fctJob, nil
}

func generateRequestHash(chainUUID string, params interface{}) (string, error) {
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	hash := md5.Sum([]byte(string(jsonParams) + chainUUID))
	return hex.EncodeToString(hash[:]), nil
}
