package transactions

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
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
	validator          validators.TransactionValidator
	txRequestDataAgent store.TransactionRequestAgent
	startJobUsecase    jobs.StartJobUseCase
}

// NewSendTxUseCase creates a new SendTxUseCase
func NewSendTxUseCase(validator validators.TransactionValidator, txRequestDataAgent store.TransactionRequestAgent, startJobUsecase jobs.StartJobUseCase) SendTxUseCase {
	return &sendTxUsecase{
		validator:          validator,
		txRequestDataAgent: txRequestDataAgent,
		startJobUsecase:    startJobUsecase,
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

	// TODO: Craft "Data" field here with MethodSignature and Args
	txData := ""
	jsonParams, err := utils.ObjectToJSON(txRequest.Params)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Create model and insert in DB
	txRequestModel := &models.TransactionRequest{
		IdempotencyKey: txRequest.IdempotencyKey,
		Schedule: &models.Schedule{
			TenantID: tenantID,
			ChainID:  txRequest.ChainID,
			Jobs: []*models.Job{{
				Type: types.JobConstantinopleTransaction,
				Transaction: &models.Transaction{
					Sender:    txRequest.Params.To,
					Recipient: txRequest.Params.From,
					Value:     txRequest.Params.Value,
					GasPrice:  txRequest.Params.GasPrice,
					GasLimit:  txRequest.Params.Gas,
					Data:      txData,
				},
				Logs: []*models.Log{{
					Status:  types.LogStatusCreated,
					Message: "Job created for contract transaction request",
				}},
				Labels: txRequest.Labels,
			}},
		},
		RequestHash: requestHash,
		Params:      jsonParams,
	}
	err = uc.txRequestDataAgent.SelectOrInsert(ctx, txRequestModel)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// Start first job of the schedule
	err = uc.startJobUsecase.Execute(ctx, txRequestModel.Schedule.Jobs[0].UUID)
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
