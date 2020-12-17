package validators

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/ethereum/abi"
	abi2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/proto"
)

//go:generate mockgen -source=transactions.go -destination=mocks/transactions.go -package=mocks

const validatorComponent = "validator"

type TransactionValidator interface {
	ValidateChainExists(ctx context.Context, chainUUID string) (string, error)
	ValidateMethodSignature(methodSignature string, args []interface{}) (string, error)
	ValidateContract(ctx context.Context, params *entities.ETHTransactionParams) (string, error)
}

// transactionValidator is a validator for transaction requests (business logic)
type transactionValidator struct {
	chainRegistryClient    client.ChainRegistryClient
	contractRegistryClient contractregistry.ContractRegistryClient
}

// NewTransactionValidator creates a new TransactionValidator
func NewTransactionValidator(
	chainRegistryClient client.ChainRegistryClient,
	contractRegistryClient contractregistry.ContractRegistryClient,
) TransactionValidator {
	return &transactionValidator{
		chainRegistryClient:    chainRegistryClient,
		contractRegistryClient: contractRegistryClient,
	}
}

func (txValidator *transactionValidator) ValidateChainExists(ctx context.Context, chainUUID string) (string, error) {
	// Validate that the chainUUID exists
	chain, err := txValidator.chainRegistryClient.GetChainByUUID(ctx, chainUUID)
	if err == nil {
		return chain.ChainID, nil
	}

	if errors.IsNotFoundError(err) {
		errMessage := "failed to get chain"
		log.WithError(err).WithField("chain_uuid", chainUUID).Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}

	log.WithError(err).WithField("chain_uuid", chainUUID).Error("failed to validate chain")
	return "", errors.FromError(err).ExtendComponent(validatorComponent)
}

func (txValidator *transactionValidator) ValidateMethodSignature(method string, args []interface{}) (string, error) {
	crafter := abi.BaseCrafter{}
	sArgs, err := utils.ParseIArrayToStringArray(args)
	if err != nil {
		errMessage := "failed to parse method arguments"
		log.WithError(err).
			WithField("method", method).
			WithField("args", args).
			Error(errMessage)
		return "", errors.DataCorruptedError(errMessage).ExtendComponent(validatorComponent)
	}

	txDataBytes, err := crafter.CraftCall(method, sArgs...)

	if err != nil {
		errMessage := "invalid method signature"
		log.WithError(err).
			WithField("method", method).
			WithField("args", args).
			Error(errMessage)

		return "", errors.InvalidParameterError(errMessage)
	}

	return hexutil.Encode(txDataBytes), nil
}

func (txValidator *transactionValidator) ValidateContract(ctx context.Context, params *entities.ETHTransactionParams) (string, error) {
	logger := log.WithContext(ctx).WithField("contract_name", params.ContractName).WithField("contract_tag", params.ContractTag)
	logger.Debug("validating contract")

	if params.ContractTag == "" {
		params.ContractTag = "latest"
	}

	response, err := txValidator.contractRegistryClient.GetContract(ctx, &contractregistry.GetContractRequest{
		ContractId: &abi2.ContractId{
			Name: params.ContractName,
			Tag:  params.ContractTag,
		},
	})
	if err != nil {
		errMessage := "failed to fetch contract"
		logger.Error(errMessage)
		return "", errors.InvalidParameterError(errMessage).ExtendComponent(validatorComponent)
	}

	// Craft bytecode
	bytecode, err := hexutil.Decode(response.Contract.GetBytecode())
	if err != nil {
		errMessage := "failed to decode bytecode"
		logger.WithError(err).Error(errMessage)
		return "", errors.DataCorruptedError(errMessage).ExtendComponent(validatorComponent)
	}

	// Craft constructor method signature
	constructorSignature := fmt.Sprintf("constructor%s", response.Contract.Constructor.Signature)
	crafter := abi.BaseCrafter{}
	args, err := utils.ParseIArrayToStringArray(params.Args)
	if err != nil {
		errMessage := "failed to parse constructor method arguments"
		logger.WithError(err).Error(errMessage)
		return "", errors.DataCorruptedError(errMessage).ExtendComponent(validatorComponent)
	}

	txDataBytes, err := crafter.CraftConstructor(bytecode, constructorSignature, args...)
	if err != nil {
		errMessage := "invalid arguments for constructor method signature"
		log.WithError(err).Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}

	return hexutil.Encode(txDataBytes), nil
}
