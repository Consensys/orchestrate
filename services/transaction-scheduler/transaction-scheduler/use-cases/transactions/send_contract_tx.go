package transactions

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"
)

//go:generate mockgen -source=send_contract_tx.go -destination=mocks/send_contract_tx.go -package=mocks

const sendContractTxComponent = "use-cases.send-contract-tx"

// sendTxUsecase is a use case to create a new contract transaction
type sendContractTxUsecase struct {
	validator     validators.TransactionValidator
	sendTxUseCase usecases.SendTxUseCase
}

// NewSendContractTxUseCase creates a n¬ew SendContractTxUseCase
func NewSendContractTxUseCase(validator validators.TransactionValidator, sendTxUseCase usecases.SendTxUseCase) usecases.SendContractTxUseCase {
	return &sendContractTxUsecase{
		validator:     validator,
		sendTxUseCase: sendTxUseCase,
	}
}

// Execute validates, creates and starts a new contract transaction
func (uc *sendContractTxUsecase) Execute(ctx context.Context, txRequest *entities.TxRequest, tenantID string) (*entities.TxRequest, error) {
	logger := log.WithContext(ctx)
	logger.WithField("idempotency_key", txRequest.IdempotencyKey).
		Debug("creating new contract transaction")

	txData, err := uc.validator.ValidateMethodSignature(txRequest.Params.MethodSignature, txRequest.Params.Args)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendContractTxComponent)
	}

	return uc.sendTxUseCase.Execute(ctx, txRequest, txData, tenantID)
}
