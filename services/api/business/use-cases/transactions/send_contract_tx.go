package transactions

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethereum/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

const sendContractTxComponent = "use-cases.send-contract-tx"

// sendTxUsecase is a use case to create a new contract transaction
type sendContractTxUsecase struct {
	sendTxUseCase usecases.SendTxUseCase
}

// NewSendContractTxUseCase creates a n¬ew SendContractTxUseCase
func NewSendContractTxUseCase(sendTxUseCase usecases.SendTxUseCase) usecases.SendContractTxUseCase {
	return &sendContractTxUsecase{
		sendTxUseCase: sendTxUseCase,
	}
}

// Execute validates, creates and starts a new contract transaction
func (uc *sendContractTxUsecase) Execute(ctx context.Context, txRequest *entities.TxRequest, tenantID string) (*entities.TxRequest, error) {
	logger := log.WithContext(ctx)
	logger.WithField("idempotency_key", txRequest.IdempotencyKey).
		Debug("creating new contract transaction")

	txData, err := uc.computeTxData(txRequest.Params.MethodSignature, txRequest.Params.Args)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendContractTxComponent)
	}

	return uc.sendTxUseCase.Execute(ctx, txRequest, txData, tenantID)
}

func (uc *sendContractTxUsecase) computeTxData(method string, args []interface{}) (string, error) {
	crafter := abi.BaseCrafter{}
	sArgs, err := utils.ParseIArrayToStringArray(args)
	if err != nil {
		errMessage := "failed to parse method arguments"
		log.WithError(err).
			WithField("method", method).
			WithField("args", args).
			Error(errMessage)
		return "", errors.DataCorruptedError(errMessage).ExtendComponent(sendContractTxComponent)
	}

	txDataBytes, err := crafter.CraftCall(method, sArgs...)

	if err != nil {
		errMessage := "invalid method signature"
		log.WithError(err).
			WithField("method", method).
			WithField("args", args).
			Error(errMessage)

		return "", errors.InvalidParameterError(errMessage)
	}

	return hexutil.Encode(txDataBytes), nil
}
