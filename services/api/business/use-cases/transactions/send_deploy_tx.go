package transactions

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethereum/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
)

const sendDeployTxComponent = "use-cases.send-deploy-tx"

// sendDeployTxUsecase is a use case to create a new contract deployment transaction
type sendDeployTxUsecase struct {
	sendTxUseCase      usecases.SendTxUseCase
	getContractUseCase usecases.GetContractUseCase
}

func NewSendDeployTxUseCase(sendTxUC usecases.SendTxUseCase, getContractUC usecases.GetContractUseCase) usecases.SendDeployTxUseCase {
	return &sendDeployTxUsecase{
		getContractUseCase: getContractUC,
		sendTxUseCase:      sendTxUC,
	}
}

// Execute validates, creates and starts a new contract deployment transaction
func (uc *sendDeployTxUsecase) Execute(ctx context.Context, txRequest *entities.TxRequest, tenantID string) (*entities.TxRequest, error) {
	logger := log.WithContext(ctx)
	logger.WithField("idempotency_key", txRequest.IdempotencyKey).Debug("creating new deployment transaction")

	txData, err := uc.computeTxData(ctx, txRequest.Params)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendDeployTxComponent)
	}

	return uc.sendTxUseCase.Execute(ctx, txRequest, txData, tenantID)
}

func (uc *sendDeployTxUsecase) computeTxData(ctx context.Context, params *entities.ETHTransactionParams) (string, error) {
	logger := log.WithContext(ctx).WithField("contract_name", params.ContractName).WithField("contract_tag", params.ContractTag)
	logger.Debug("validating contract")

	if params.ContractTag == "" {
		params.ContractTag = "latest"
	}

	contract, err := uc.getContractUseCase.Execute(ctx, &entities.ContractID{
		Name: params.ContractName,
		Tag:  params.ContractTag,
	})

	if err != nil {
		errMessage := "failed to fetch contract"
		logger.Error(errMessage)
		return "", errors.InvalidParameterError(errMessage).ExtendComponent(sendDeployTxComponent)
	}

	// Craft bytecode
	bytecode, err := hexutil.Decode(contract.Bytecode)
	if err != nil {
		errMessage := "failed to decode bytecode"
		logger.WithError(err).Error(errMessage)
		return "", errors.DataCorruptedError(errMessage).ExtendComponent(sendDeployTxComponent)
	}

	// Craft constructor method signature
	constructorSignature := fmt.Sprintf("constructor%s", contract.Constructor.Signature)
	crafter := abi.BaseCrafter{}
	args, err := utils.ParseIArrayToStringArray(params.Args)
	if err != nil {
		errMessage := "failed to parse constructor method arguments"
		logger.WithError(err).Error(errMessage)
		return "", errors.DataCorruptedError(errMessage).ExtendComponent(sendDeployTxComponent)
	}

	txDataBytes, err := crafter.CraftConstructor(bytecode, constructorSignature, args...)
	if err != nil {
		errMessage := "invalid arguments for constructor method signature"
		log.WithError(err).Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}

	return hexutil.Encode(txDataBytes), nil
}
