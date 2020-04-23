package validators

import (
	"context"

	log "github.com/sirupsen/logrus"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
)

//go:generate mockgen -source=transactions.go -destination=mocks/transactions.go -package=mocks

const txValidatorComponent = "transaction-validator"

type TransactionValidator interface {
	ValidateTx(ctx context.Context, txRequest *types.TransactionRequest) error
}

// Transaction is a validator for transaction requests (business logic)
type Transaction struct {
	chainRegistryClient client.ChainRegistryClient
}

// NewTransaction creates a new TransactionValidator
func NewTransaction(chainRegistryClient client.ChainRegistryClient) TransactionValidator {
	return &Transaction{chainRegistryClient: chainRegistryClient}
}

// ValidateTx validates a transaction request
func (validator *Transaction) ValidateTx(ctx context.Context, txRequest *types.TransactionRequest) error {
	log.WithContext(ctx).WithField("tx", txRequest).Debug("validating transaction")

	// TODO: Validation of args given methodSignature

	_, err := validator.chainRegistryClient.GetChainByUUID(ctx, txRequest.ChainID)
	if err != nil {
		return errors.FromError(err).ExtendComponent(txValidatorComponent)
	}

	return nil
}
