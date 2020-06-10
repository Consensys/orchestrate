package usecases

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/abi"
)

const (
	getMethodSignaturesComponent = component + ".get-method-signatures"
	constructorMethodName        = "constructor"
)

//go:generate mockgen -source=get_method_signatures.go -destination=mocks/mock_get_method_signatures.go -package=mocks

type GetMethodSignaturesUseCase interface {
	Execute(ctx context.Context, contract *abi.ContractId, methodName string) ([]string, error)
}

// GetMethods is a use case to get methods
type GetMethodSignatures struct {
	getContractUseCase GetContractUseCase
}

// NewGetMethodSignatures creates a new GetMethodSignatures
func NewGetMethodSignatures(getContractUseCase GetContractUseCase) *GetMethodSignatures {
	return &GetMethodSignatures{
		getContractUseCase: getContractUseCase,
	}
}

// Execute validates and registers a new contract in DB
func (uc *GetMethodSignatures) Execute(ctx context.Context, id *abi.ContractId, methodName string) ([]string, error) {
	logger := log.WithContext(ctx).WithField("contract", id).WithField("method_name", methodName)
	logger.Debug("getting method signatures")

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
			if method.Name == methodName {
				signatures = append(signatures, method.Sig())
			}
		}
	}

	return signatures, nil
}
