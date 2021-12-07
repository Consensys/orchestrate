package transactions

import (
	"context"
	"fmt"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/ethereum/abi"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/orchestrate/pkg/utils"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
)

const sendDeployTxComponent = "use-cases.send-deploy-tx"

// sendDeployTxUsecase is a use case to create a new contract deployment transaction
type sendDeployTxUsecase struct {
	sendTxUseCase      usecases.SendTxUseCase
	getContractUseCase usecases.GetContractUseCase
	logger             *log.Logger
}

func NewSendDeployTxUseCase(sendTxUC usecases.SendTxUseCase, getContractUC usecases.GetContractUseCase) usecases.SendDeployTxUseCase {
	return &sendDeployTxUsecase{
		getContractUseCase: getContractUC,
		sendTxUseCase:      sendTxUC,
		logger:             log.NewLogger().SetComponent(sendDeployTxComponent),
	}
}

// Execute validates, creates and starts a new contract deployment transaction
func (uc *sendDeployTxUsecase) Execute(ctx context.Context, txRequest *entities.TxRequest, userInfo *multitenancy.UserInfo) (*entities.TxRequest, error) {
	ctx = log.WithFields(ctx, log.Field("idempotency-key", txRequest.IdempotencyKey))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("creating new deployment transaction")

	txData, err := uc.computeTxData(ctx, txRequest.Params)
	if err != nil {
		logger.WithError(err).WithField("contract_name", txRequest.Params.ContractName).
			WithField("contract_tag", txRequest.Params.ContractTag).Error("failed to compute transaction data")
		return nil, errors.FromError(err).ExtendComponent(sendDeployTxComponent)
	}

	return uc.sendTxUseCase.Execute(ctx, txRequest, txData, userInfo)
}

func (uc *sendDeployTxUsecase) computeTxData(ctx context.Context, params *entities.ETHTransactionParams) ([]byte, error) {
	logger := uc.logger.WithContext(ctx)

	if params.ContractTag == "" {
		params.ContractTag = "latest"
	}

	contract, err := uc.getContractUseCase.Execute(ctx, params.ContractName, params.ContractTag)
	if errors.IsNotFoundError(err) {
		return nil, errors.InvalidParameterError("contract not found")
	}
	if err != nil {
		return nil, errors.FromError(err)
	}

	// Craft constructor method signature
	constructorSignature := fmt.Sprintf("constructor%s", contract.Constructor.Signature)
	crafter := abi.BaseCrafter{}
	args, err := utils.ParseIArrayToStringArray(params.Args)
	if err != nil {
		errMessage := "failed to parse constructor method arguments"
		logger.WithError(err).WithField("args", params.Args).Error(errMessage)
		return nil, errors.DataCorruptedError(errMessage)
	}

	txData, err := crafter.CraftConstructor(contract.Bytecode, constructorSignature, args...)
	if err != nil {
		errMessage := "invalid arguments for constructor method signature"
		logger.WithError(err).WithField("signature", constructorSignature).WithField("args", args).Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}

	return txData, nil
}
