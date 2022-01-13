package transactions

import (
	"context"

	"github.com/umbracle/go-web3/abi"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
	"github.com/consensys/orchestrate/pkg/types/entities"
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
	ctx = log.WithFields(
		ctx,
		log.Field("idempotency-key", txRequest.IdempotencyKey),
		log.Field("method", txRequest.Params.MethodSignature),
		log.Field("args", txRequest.Params.Args),
	)
	logger := uc.logger.WithContext(ctx)
	logger.Debug("creating new deployment transaction")

	contract, err := uc.getContractUseCase.Execute(ctx, txRequest.Params.ContractName, txRequest.Params.ContractTag)
	if errors.IsNotFoundError(err) {
		return nil, errors.InvalidParameterError("contract not found")
	}
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(sendDeployTxComponent)
	}

	if contract.Bytecode == nil || len(contract.Bytecode) == 0 {
		errMessage := "contract has no bytecode"
		uc.logger.WithError(err).Error(errMessage)
		return nil, errors.DataCorruptedError(errMessage).ExtendComponent(sendDeployTxComponent)
	}

	// TODO: We restrict the usage of web3-go to only generate the txData but ideally we should use it as much as possible and change the ABI type everywhere in the codebase
	web3ABI, err := abi.NewABI(contract.RawABI)
	if err != nil {
		errMessage := "failed to parse contract ABI"
		logger.WithError(err).Error(errMessage)
		return nil, errors.DataCorruptedError(errMessage).ExtendComponent(sendContractTxComponent)
	}

	var arguments []byte
	if web3ABI.Constructor != nil { // It is possible to create a smart contract without constructor
		arguments, err = abi.Encode(txRequest.Params.Args, web3ABI.Constructor.Inputs)
		if err != nil {
			logger.WithError(err).Error("failed to compute tx data from constructor and arguments")
			return nil, errors.InvalidParameterError(err.Error()).ExtendComponent(sendContractTxComponent)
		}
	}

	txData := append(contract.Bytecode, arguments...)
	return uc.sendTxUseCase.Execute(ctx, txRequest, txData, userInfo)
}
