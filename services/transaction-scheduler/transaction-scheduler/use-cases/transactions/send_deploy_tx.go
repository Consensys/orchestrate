package transactions

import (
	"context"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/use-cases"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/validators"
)

//go:generate mockgen -source=send_deploy_tx.go -destination=mocks/send_deploy_tx.go -package=mocks

const sendDeployTxComponent = "use-cases.send-deploy-tx"

// sendDeployTxUsecase is a use case to create a new contract deployment transaction
type sendDeployTxUsecase struct {
	validator     validators.TransactionValidator
	sendTxUseCase usecases.SendTxUseCase
}

// NewSendDeployTxUseCase creates a new SendDeployTxUseCase
func NewSendDeployTxUseCase(validator validators.TransactionValidator, sendTxUseCase usecases.SendTxUseCase) usecases.SendDeployTxUseCase {
	return &sendDeployTxUsecase{
		validator:     validator,
		sendTxUseCase: sendTxUseCase,
	}
}

// Execute validates, creates and starts a new contract deployment transaction
func (uc *sendDeployTxUsecase) Execute(ctx context.Context, txRequest *entities.TxRequest, tenantID string) (*entities.TxRequest, error) {
	logger := log.WithContext(ctx)
	logger.WithField("idempotency_key", txRequest.IdempotencyKey).Debug("creating new deployment transaction")

	txData, err := uc.validator.ValidateContract(ctx, txRequest.Params)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendDeployTxComponent)
	}

	return uc.sendTxUseCase.Execute(ctx, txRequest, txData, tenantID)
}
