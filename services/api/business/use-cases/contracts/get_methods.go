package contracts

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/store"
)

const getMethodsComponent = "use-cases.get-methods"

type getMethodsUseCase struct {
	agent store.MethodAgent
}

func NewGetMethodsUseCase(agent store.MethodAgent) usecases.GetContractMethodsUseCase {
	return &getMethodsUseCase{
		agent: agent,
	}
}

func (usecase *getMethodsUseCase) Execute(ctx context.Context, chainID, address string, selector []byte) (abi string, methodsABI []string, err error) {
	logger := log.WithContext(ctx).WithField("chainID", chainID).WithField("address", address)
	logger.Debug("get methods starting...")

	method, err := usecase.agent.FindOneByAccountAndSelector(ctx, chainID, address, selector)
	if errors.IsConnectionError(err) {
		return "", nil, errors.FromError(err).ExtendComponent(getMethodsComponent)
	}
	if method != nil {
		return method.ABI, nil, nil
	}

	defaultMethods, err := usecase.agent.FindDefaultBySelector(ctx, selector)
	if err != nil {
		return "", nil, errors.FromError(err).ExtendComponent(getMethodsComponent)
	}

	for _, m := range defaultMethods {
		methodsABI = append(methodsABI, m.ABI)
	}

	logger.Debug("get methods executed successfully")
	return "", methodsABI, nil
}
