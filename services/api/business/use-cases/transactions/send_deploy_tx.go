package transactions

import (
	"context"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/ethereum/abi"
	"github.com/ConsenSys/orchestrate/pkg/types/entities"
	"github.com/ConsenSys/orchestrate/pkg/utils"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
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
func (uc *sendDeployTxUsecase) Execute(ctx context.Context, txRequest *entities.TxRequest, tenantID string) (*entities.TxRequest, error) {
	ctx = log.WithFields(ctx, log.Field("idempotency-key", txRequest.IdempotencyKey))
	logger := uc.logger.WithContext(ctx)
	logger.Debug("creating new deployment transaction")

	txData, err := uc.computeTxData(ctx, txRequest.Params)
	if err != nil {
		logger.WithError(err).WithField("contract_name", txRequest.Params.ContractName).
			WithField("contract_tag", txRequest.Params.ContractTag).Error("failed to compute transaction data")
		return nil, errors.FromError(err).ExtendComponent(sendDeployTxComponent)
	}

	return uc.sendTxUseCase.Execute(ctx, txRequest, txData, tenantID)
}

func (uc *sendDeployTxUsecase) computeTxData(ctx context.Context, params *entities.ETHTransactionParams) (string, error) {
	if params.ContractTag == "" {
		params.ContractTag = "latest"
	}

	contract, err := uc.getContractUseCase.Execute(ctx, params.ContractName, params.ContractTag)
	if errors.IsNotFoundError(err) {
		return "", errors.InvalidParameterError("contract not found")
	}
	if err != nil {
		return "", errors.FromError(err)
	}

	// Craft bytecode
	bytecode, err := hexutil.Decode(contract.Bytecode)
	if err != nil {
		return "", errors.EncodingError("failed to decode bytecode")
	}

	// Craft constructor method signature
	constructorSignature := fmt.Sprintf("constructor%s", contract.Constructor.Signature)
	crafter := abi.BaseCrafter{}
	args, err := utils.ParseIArrayToStringArray(params.Args)
	if err != nil {
		return "", errors.DataCorruptedError("failed to parse constructor method arguments")
	}

	txDataBytes, err := crafter.CraftConstructor(bytecode, constructorSignature, args...)
	if err != nil {
		return "", errors.InvalidParameterError("invalid arguments for constructor method signature")
	}

	return hexutil.Encode(txDataBytes), nil
}
