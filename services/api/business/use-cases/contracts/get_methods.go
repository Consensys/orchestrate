package contracts

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	usecases "github.com/consensys/orchestrate/services/api/business/use-cases"
	"github.com/consensys/orchestrate/services/api/store"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const getMethodsComponent = "use-cases.get-methods"

type getMethodsUseCase struct {
	agent  store.MethodAgent
	logger *log.Logger
}

func NewGetMethodsUseCase(agent store.MethodAgent) usecases.GetContractMethodsUseCase {
	return &getMethodsUseCase{
		agent:  agent,
		logger: log.NewLogger().SetComponent(getMethodsComponent),
	}
}

func (uc *getMethodsUseCase) Execute(ctx context.Context, chainID string, address ethcommon.Address, selector []byte) (abi string, methodsABI []string, err error) {
	ctx = log.WithFields(ctx, log.Field("chain_id", chainID), log.Field("address", address))
	logger := uc.logger.WithContext(ctx)

	method, err := uc.agent.FindOneByAccountAndSelector(ctx, chainID, address.Hex(), selector)
	if errors.IsConnectionError(err) {
		return "", nil, errors.FromError(err).ExtendComponent(getMethodsComponent)
	}
	if method != nil {
		return method.ABI, nil, nil
	}

	defaultMethods, err := uc.agent.FindDefaultBySelector(ctx, selector)
	if err != nil {
		return "", nil, errors.FromError(err).ExtendComponent(getMethodsComponent)
	}

	for _, m := range defaultMethods {
		methodsABI = append(methodsABI, m.ABI)
	}

	logger.Debug("get methods executed successfully")
	return "", methodsABI, nil
}
