package contracts

import (
	"context"
	"fmt"

	"github.com/ConsenSys/orchestrate/pkg/errors"
	"github.com/ConsenSys/orchestrate/pkg/toolkit/app/log"
	usecases "github.com/ConsenSys/orchestrate/services/api/business/use-cases"
)

const (
	getMethodSignaturesComponent = "use-cases.get-method-signatures"
	constructorMethodName        = "constructor"
)

type getMethodSignaturesUseCase struct {
	getContractUseCase usecases.GetContractUseCase
	logger             *log.Logger
}

func NewGetMethodSignaturesUseCase(getContractUseCase usecases.GetContractUseCase) usecases.GetContractMethodSignaturesUseCase {
	return &getMethodSignaturesUseCase{
		getContractUseCase: getContractUseCase,
		logger:             log.NewLogger().SetComponent(getMethodSignaturesComponent),
	}
}

func (uc *getMethodSignaturesUseCase) Execute(ctx context.Context, name, tag, methodName string) ([]string, error) {
	ctx = log.WithFields(ctx, log.Field("contract_name", name), log.Field("tag", tag))
	logger := uc.logger.WithContext(ctx)

	contract, err := uc.getContractUseCase.Execute(ctx, name, tag)
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

	logger.Debug("contract method signatures were fetched successfully")
	return signatures, nil
}
