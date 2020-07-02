package validators

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/ethereum/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types"
	abi2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/abi"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	contractregistry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/proto"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/store"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/entities"
)

//go:generate mockgen -source=transactions.go -destination=mocks/transactions.go -package=mocks

const txValidatorComponent = "transaction-validator"

type TransactionValidator interface {
	ValidateFields(ctx context.Context, txRequest *entities.TxRequest) error
	ValidateChainExists(ctx context.Context, chainUUID string) error
	ValidateMethodSignature(methodSignature string, args []interface{}) (string, error)
	ValidateContract(ctx context.Context, params *types.ETHTransactionParams) (string, error)
}

// transactionValidator is a validator for transaction requests (business logic)
type transactionValidator struct {
	db                     store.DB
	chainRegistryClient    client.ChainRegistryClient
	contractRegistryClient contractregistry.ContractRegistryClient
}

// NewTransactionValidator creates a new TransactionValidator
func NewTransactionValidator(
	db store.DB,
	chainRegistryClient client.ChainRegistryClient,
	contractRegistryClient contractregistry.ContractRegistryClient,
) TransactionValidator {
	return &transactionValidator{db: db, chainRegistryClient: chainRegistryClient, contractRegistryClient: contractRegistryClient}
}

func (txValidator *transactionValidator) ValidateFields(ctx context.Context, txRequest *entities.TxRequest) error {
	logger := log.WithContext(ctx)

	if err := utils.GetValidator().Struct(txRequest); err != nil {
		errMessage := err.Error()
		logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(errMessage).ExtendComponent(txValidatorComponent)
	}

	if err := txRequest.Params.PrivateTransactionParams.Validate(); err != nil {
		errMessage := err.Error()
		logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(err.Error()).ExtendComponent(txValidatorComponent)
	}

	if txRequest.Annotations != nil && txRequest.Annotations.OneTimeKey && txRequest.Params.From != "" {
		errMessage := "from account cannot be included when OneTimeKey is enabled"
		logger.Error(errMessage)
		return errors.InvalidParameterError(errMessage).ExtendComponent(txValidatorComponent)
	}

	return nil
}

func (txValidator *transactionValidator) ValidateChainExists(ctx context.Context, chainUUID string) error {
	// Validate that the chainUUID exists
	_, err := txValidator.chainRegistryClient.GetChainByUUID(ctx, chainUUID)
	if err != nil {
		errMessage := "failed to get chain"
		log.WithError(err).WithField("chain_uuid", chainUUID).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	return nil
}

func (txValidator *transactionValidator) ValidateMethodSignature(method string, args []interface{}) (string, error) {
	crafter := abi.BaseCrafter{}
	txDataBytes, err := crafter.CraftCall(method, utils.ParseIArrayToStringArray(args)...)

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

func (txValidator *transactionValidator) ValidateContract(ctx context.Context, params *types.ETHTransactionParams) (string, error) {
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
		return "", errors.InvalidParameterError(errMessage).ExtendComponent(txValidatorComponent)
	}

	// Craft bytecode
	bytecode, err := hexutil.Decode(response.Contract.GetBytecode())
	if err != nil {
		errMessage := "failed to decode bytecode"
		logger.WithError(err).Error(errMessage)
		return "", errors.DataCorruptedError(errMessage).ExtendComponent(txValidatorComponent)
	}

	// Craft constructor method signature
	constructorSignature := fmt.Sprintf("constructor%s", response.Contract.Constructor.Signature)
	crafter := abi.BaseCrafter{}
	txDataBytes, err := crafter.CraftConstructor(bytecode, constructorSignature, utils.ParseIArrayToStringArray(params.Args)...)
	if err != nil {
		errMessage := "invalid arguments for constructor method signature"
		log.WithError(err).Error(errMessage)
		return "", errors.InvalidParameterError(errMessage)
	}

	return hexutil.Encode(txDataBytes), nil
}
