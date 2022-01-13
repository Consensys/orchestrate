package transactions

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/umbracle/go-web3/abi"
)

const sendContractTxComponent = "use-cases.send-contract-tx"

type sendContractTxUseCase struct {
	sendTxUseCase      usecases.SendTxUseCase
	getContractUseCase usecases.GetContractUseCase
	logger             *log.Logger
}

// NewSendContractTxUseCase creates a n¬ew SendContractTxUseCase
func NewSendContractTxUseCase(sendTxUseCase usecases.SendTxUseCase, getContractUseCase usecases.GetContractUseCase) usecases.SendContractTxUseCase {
	return &sendContractTxUseCase{
		sendTxUseCase:      sendTxUseCase,
		getContractUseCase: getContractUseCase,
		logger:             log.NewLogger().SetComponent(sendContractTxComponent),
	}
}

// Execute validates, creates and starts a new contract transaction
func (uc *sendContractTxUseCase) Execute(ctx context.Context, txRequest *entities.TxRequest, userInfo *multitenancy.UserInfo) (*entities.TxRequest, error) {
	ctx = log.WithFields(
		ctx,
		log.Field("idempotency-key", txRequest.IdempotencyKey),
		log.Field("method", txRequest.Params.MethodSignature),
		log.Field("args", txRequest.Params.Args),
	)
	logger := uc.logger.WithContext(ctx)
	logger.Debug("creating new contract transaction")

	contract, err := uc.getContractUseCase.Execute(ctx, txRequest.Params.ContractName, txRequest.Params.ContractTag)
	if errors.IsNotFoundError(err) {
		return nil, errors.InvalidParameterError("contract not found")
	}
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendContractTxComponent)
	}

	// TODO: We restrict the usage of web3-go to only generate the txData but ideally we should use it as much as possible and change the ABI type everywhere in the codebase
	web3ABI, err := abi.NewABI(contract.RawABI)
	if err != nil {
		errMessage := "failed to parse contract ABI for contract transaction"
		logger.WithError(err).Error(errMessage)
		return nil, errors.DataCorruptedError(errMessage).ExtendComponent(sendContractTxComponent)
	}

	method := web3ABI.GetMethodBySignature(txRequest.Params.MethodSignature)
	if method == nil {
		errMessage := "method not found"
		logger.WithError(err).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage).ExtendComponent(sendContractTxComponent)
	}

	txData, err := method.Encode(txRequest.Params.Args)
	if err != nil {
		logger.WithError(err).Error("failed to compute tx data from method signature and arguments")
		return nil, errors.InvalidParameterError(err.Error()).ExtendComponent(sendContractTxComponent)
	}

	tx, err := uc.sendTxUseCase.Execute(ctx, txRequest, txData, userInfo)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendContractTxComponent)
	}

	return tx, nil
}
