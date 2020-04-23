package transactions

import (
	"context"

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
	txRequestDataAgent store.TransactionRequestAgent
	validator          validators.TransactionValidator
}

// NewSendTx creates a new SendTxUseCase
func NewSendTx(txRequestDataAgent store.TransactionRequestAgent, validator validators.TransactionValidator) SendTxUseCase {
	return &SendTx{
		txRequestDataAgent: txRequestDataAgent,
		validator:          validator,
	}
}

// Execute validates, creates and starts a new transaction
func (usecase *SendTx) Execute(ctx context.Context, txRequest *types.TransactionRequest, tenantID string) (*types.TransactionResponse, error) {
	log.WithContext(ctx).WithField("idempotency_key", txRequest.IdempotencyKey).Info("creating new transaction")

	err := usecase.validator.ValidateTx(ctx, txRequest)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	jsonParams, err := utils.ObjectToJSON(txRequest.Params)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	txRequestModel := &models.TransactionRequest{
		IdempotencyKey: txRequest.IdempotencyKey,
		Chain:          txRequest.ChainID,
		Method:         types.MethodSendRawTransaction,
		Params:         jsonParams,
		Labels:         txRequest.Labels,
	}
	err = usecase.txRequestDataAgent.Insert(ctx, txRequestModel)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendTxComponent)
	}

	// TODO: Create schedule (with use case) that also creates jobs

	// TODO: Start job

	return utils.FormatTxResponse(ctx, txRequestModel)
}
