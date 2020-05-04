package transactions

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/jobs"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases/schedules"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"
)

//go:generate mockgen -source=send_tx.go -destination=mocks/send_tx.go -package=mocks

const sendTxComponent = "use-cases.send-tx"

type SendTxUseCase interface {
	Execute(ctx context.Context, txRequest *types.TransactionRequest, tenantID string) (*types.TransactionResponse, error)
}

// SendTx is a use case to create a new transaction request
type SendTx struct {
	createScheduleUseCase schedules.CreateScheduleUseCase
	createJobUseCase      jobs.CreateJobUseCase
	txRequestDataAgent    store.TransactionRequestAgent
	validator             validators.TransactionValidator
}

// NewSendTx creates a new SendTxUseCase
func NewSendTx(
	txRequestDataAgent store.TransactionRequestAgent,
	validator validators.TransactionValidator,
	createScheduleUseCase schedules.CreateScheduleUseCase,
	createJobUseCase jobs.CreateJobUseCase,
) SendTxUseCase {
	return &SendTx{
		createScheduleUseCase: createScheduleUseCase,
		createJobUseCase:      createJobUseCase,
		txRequestDataAgent:    txRequestDataAgent,
		validator:             validator,
	}
}

// Execute validates, creates and starts a new transaction
func (usecase *SendTx) Execute(ctx context.Context, txRequest *types.TransactionRequest, tenantID string) (*types.TransactionResponse, error) {
	log.WithContext(ctx).WithField("idempotency_key", txRequest.IdempotencyKey).Debug("creating new transaction")

	// TODO: Validation of args given methodSignature

	requestHash, err := usecase.validator.ValidateRequestHash(ctx, txRequest.Params, txRequest.IdempotencyKey)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// We create the schedule
	scheduleRequest := &types.ScheduleRequest{ChainID: txRequest.ChainID}
	scheduleResponse, scheduleID, err := usecase.createScheduleUseCase.Execute(ctx, scheduleRequest, tenantID)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// We create the job as a public ETH transaction
	// TODO: Craft "Data" field here with MethodSignature and Args
	txData := ""
	jobRequest := &types.JobRequest{
		ScheduleID: scheduleID,
		Type:       types.JobConstantinopleTransaction,
		Labels:     txRequest.Labels,
		Transaction: types.ETHTransaction{
			From:     &txRequest.Params.From,
			To:       &txRequest.Params.To,
			Value:    txRequest.Params.Value,
			GasPrice: txRequest.Params.GasPrice,
			GasLimit: txRequest.Params.Gas,
			Data:     &txData,
		},
	}
	jobResponse, err := usecase.createJobUseCase.Execute(ctx, jobRequest)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	jsonParams, err := utils.ObjectToJSON(txRequest.Params)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	txRequestModel := &models.TransactionRequest{
		IdempotencyKey: txRequest.IdempotencyKey,
		ScheduleID:     scheduleID,
		RequestHash:    requestHash,
		Params:         jsonParams,
	}
	err = usecase.txRequestDataAgent.SelectOrInsert(ctx, txRequestModel)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// TODO: Start necessary job and update the logs to STARTED

	scheduleResponse.Jobs = append([]*types.JobResponse{}, jobResponse)
	txResponse, err := utils.FormatTxResponse(txRequestModel, scheduleResponse)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	log.WithContext(ctx).WithField("idempotency_key", txRequestModel.IdempotencyKey).Info("contract transaction request created successfully")
	return txResponse, nil
}
