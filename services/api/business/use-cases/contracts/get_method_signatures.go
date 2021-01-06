package contracts

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/api/business/use-cases"
)

const (
	getMethodSignaturesComponent = "use-cases.get-method-signatures"
	constructorMethodName        = "constructor"
)

type getMethodSignaturesUseCase struct {
	getContractUseCase usecases.GetContractUseCase
}

func NewGetMethodSignaturesUseCase(getContractUseCase usecases.GetContractUseCase) usecases.GetContractMethodSignaturesUseCase {
	return &getMethodSignaturesUseCase{
		getContractUseCase: getContractUseCase,
	}
}

func (uc *getMethodSignaturesUseCase) Execute(ctx context.Context, id *entities.ContractID, methodName string) ([]string, error) {
	logger := log.WithContext(ctx).WithField("contract", id).WithField("method_name", methodName)
	logger.Debug("get method signatures starting...")

	contract, err := uc.getContractUseCase.Execute(ctx, id)
	if err != nil {
		return nil, errors.FromError(err).ExtendComponent(getMethodSignaturesComponent)
	}

	contractABI, err := contract.ToABI()
	if err != nil {
		errMessage := "failed to parse contract ABI"
		logger.WithError(err).Error(errMessage)
		return nil, errors.DataCorruptedError(errMessage).ExtendComponent(getMethodSignaturesComponent)
	}

	var signatures []string

	if methodName == constructorMethodName {
		signatures = append(signatures, fmt.Sprintf("%s%s", constructorMethodName, contractABI.Constructor.Sig()))
	} else {
		for _, method := range contractABI.Methods {
			if methodName == "" || method.Name == methodName {
				signatures = append(signatures, method.Sig())
			}
		}
	}

	logger.Debug("get method signatures successfully")
	return signatures, nil
}
