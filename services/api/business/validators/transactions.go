package validators

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/chain-registry/client"
)

//go:generate mockgen -source=transactions.go -destination=mocks/transactions.go -package=mocks

const validatorComponent = "validator"

type TransactionValidator interface {
	ValidateChainExists(ctx context.Context, chainUUID string) (string, error)
}

// transactionValidator is a validator for transaction requests (business logic)
type transactionValidator struct {
	chainRegistryClient client.ChainRegistryClient
}

// NewTransactionValidator creates a new TransactionValidator
func NewTransactionValidator(
	chainRegistryClient client.ChainRegistryClient,
) TransactionValidator {
	return &transactionValidator{
		chainRegistryClient: chainRegistryClient,
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
